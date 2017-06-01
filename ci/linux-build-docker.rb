require 'pathname'

DOCKER_IMAGE_UBUNTU='ocaml/opam:ubuntu-12.04_ocaml-4.04.1'
DOCKER_IMAGE_ALPINE='ocaml/opam:alpine-3.3_ocaml-4.04.1'
DOCKER_IMAGE_CENTOS='ocaml/opam:centos-6_ocaml-4.04.1'

UBUNTU_1204_LDD_VERSION = <<-EOS
ldd (Ubuntu EGLIBC 2.15-0ubuntu10.18) 2.15
Copyright (C) 2012 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.  There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
Written by Roland McGrath and Ulrich Drepper.
EOS

CENTOS_LDD_VESION = <<-EOS
ldd (GNU libc) 2.12
Copyright (C) 2010 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.  There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
Written by Roland McGrath and Ulrich Drepper.
EOS

ALPINE_LDD_VERSION = <<-EOS
musl libc
Version 1.1.12
Dynamic Program Loader
Usage: ldd [options] [--] pathname
EOS

DEST = Pathname.new(__dir__).join('build')
HERE = Pathname.new(__dir__)

class Build
  def initialize(opts = {})
    @image = opts.fetch :image
    @dest = opts.fetch :dest
    @unison_src = opts[:unison_src] || HERE.join('unison').to_s
    @extra_commands = opts.fetch(:extra_commands, [])
  end

  attr_reader :image, :unison_src, :extra_commands

  def perform
    with_container do |container|
      # add unison source to container
      run("docker cp #{unison_src} #{container}:/unison")
      run("docker exec #{container} sudo chown -R opam /unison")

      # Uncomment to interact with the container before the build starts
      #interact(container)

      # do various extra things
      extra_commands.each do |cmd|
        run("docker exec #{container} #{cmd}")
      end


      # make unison
      #
      # note the "eval opam config env", which ensures that we're actually
      # using the OPAM env for our build.
      run(%Q(docker exec #{container} bash -c 'eval $(opam config env) && cd /unison && make'))

      # retrieve built product
      run("docker cp #{container}:/unison/src/unison #{dest}/unison")
      run("docker cp #{container}:/unison/src/unison-fsmonitor #{dest}/unison-fsmonitor")
    end
  end

  private

  def interact(container)
    puts "Stopping for interaction. Exit or press Ctrl-D to continue the build"
    puts "CONTAINER: #{container}"
    run("docker exec -it #{container} /bin/bash || true")
  end

  def run(cmd)
    puts "*** run: #{cmd}"
    result = system(cmd)
    raise "Command #{cmd.inspect} failed: #{result}" unless result
    result
  end

  def capture(cmd)
    puts "*** capture: #{cmd}"
    stdout = `#{cmd}`
    raise "Command #{cmd.inspect} failed: #{$?}" unless $?.success?
    stdout
  end

  def with_container
    name = create_container
    begin
      yield(name)
    ensure
      rm_container(name)
    end
  end

  def create_container
    capture("docker run -d -it #{image} tail -f /dev/null").strip
  end

  def rm_container(name)
    run("docker rm -f #{name}")
  end

  def dest
    res = HERE.join("build", @dest)
    res.mkpath unless res.exist?
    res
  end
end

if ARGV.length >= 2
  Build.new(image: ARGV[0], unison_src: ARGV[1]).perform
else
  # eglibc - not needed? should be ok with glibc
  #Build.new(
    #image: DOCKER_IMAGE_GLIBC,
    #unison_src: ARGV[0],
    #dest: "linux-eglibc-2.15",
  #).perform

  # musl
  Build.new(
    image: DOCKER_IMAGE_ALPINE,
    unison_src: ARGV[0],
    extra_commands: [
      # force inotify support recognized, even on a non-glibc system
      # official alpine unison package basically does the same thing, see
      # https://git.alpinelinux.org/cgit/aports/diff/community/unison/fix_inotify_check.patch?id=f5fb4fa4ed695ff1c8eda1971f0d4be46bc85864
      "sed -i -e 's/GLIBC_SUPPORT_INOTIFY 0/GLIBC_SUPPORT_INOTIFY 1/' /unison/src/fsmonitor/linux/inotify_stubs.c"
    ],
    dest: "linux-musl",
  ).perform


  # Centos 6 - glibc?
  Build.new(
    image: DOCKER_IMAGE_CENTOS,
    unison_src: ARGV[0],
    dest: "linux-glibc"
  ).perform
end
