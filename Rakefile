# TODO: learn make, use make instead
# TODO: actually track dependencies, only rebuild if needed

namespace :unison do
  desc "build linux binaries via Docker"
  task :linux do
    ruby "./ci/linux-build-docker.rb"
  end
end

desc "convert external binary assets into golang source"
task :assets do
  sh 'go-bindata -o unison/assets.go ./ci/build/...'
end

task :help do
  sh "rake -T"
end

task default: [:help]
