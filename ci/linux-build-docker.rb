require 'pathname'
DOCKER_IMAGE_GLIBC='ocaml/opam:ubuntu-12.04_ocaml-4.04.1'
DOCKER_IMAGE_MUSL='ocaml/opam:alpine-3.3_ocaml-4.04.1'
DEST = Pathname.new(__dir__).join('build')

class Build
  def initialize(opts = {})
    @image = opts[:image] || DOCKER_IMAGE_GLIBC
    @unison_src = opts[:unison_src] || "./unison"
    @extra_commands = opts[:extra_commands] || []
  end

  attr_reader :image, :unison_src, :extra_commands

  def perform
    with_container do |container|
      # add unison source to container
      run("docker cp #{unison_src} #{container}:/unison")
      run("docker exec #{container} sudo chown -R opam /unison")

      # do various extra things
      extra_commands.each do |cmd|
        run("docker exec #{container} #{cmd}")
      end

      # make unison
      run(%Q(docker exec #{container} bash -c 'cd /unison && make'))

      # retrieve built product
      run("docker cp #{container}:/unison/src/unison #{dest}/unison")
      run("docker cp #{container}:/unison/src/unison-fsmonitor #{dest}/unison-fsmonitor")
    end
  end

  private

  def interact(container)
    puts "Stopping for interaction. Exit or press Ctrl-D to continue"
    puts "CONTAINER: #{container}"
    run("docker exec -it #{container} /bin/bash")
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
    @dest ||= begin
      os = image.split(':').last
      res = Pathname.new(__dir__).join("build/#{os}")
      res.mkpath
      res
    end
  end
end

if ARGV.length >= 2
  Build.new(image: ARGV[0], unison_src: ARGV[1]).perform
else
  #Build.new(image: DOCKER_IMAGE_GLIBC, unison_src: ARGV[0]).perform
  Build.new(
    image: DOCKER_IMAGE_MUSL,
    unison_src: ARGV[0],
    extra_commands: [
      # force inotify support recognized, even on a non-glibc system
      # official alpine unison package basically does the same thing, see
      # https://git.alpinelinux.org/cgit/aports/diff/community/unison/fix_inotify_check.patch?id=f5fb4fa4ed695ff1c8eda1971f0d4be46bc85864
      "sed -i -e 's/GLIBC_SUPPORT_INOTIFY 0/GLIBC_SUPPORT_INOTIFY 1/' /unison/src/fsmonitor/linux/inotify_stubs.c"
    ],
  ).perform
end
