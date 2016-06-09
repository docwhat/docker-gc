require 'tmpdir'

PACKAGES = FileList['**/*.go']
  .exclude(/_test\.go\Z/)
  .map { |p| './' + File.dirname(p) }
  .sort
  .uniq
  .freeze

task default: %w(lint test build)

desc 'Fetches development and build dependencies'
task :setup do
  sh 'go get -u -v '\
    ' github.com/alecthomas/gometalinter'\
    ' github.com/mitchellh/gox '\
    ' golang.org/x/tools/cmd/cover'
  sh 'gometalinter --install --update'
end

desc 'Check lint and style'
task :lint do
  sh 'gometalinter --deadline=1m --disable=gotype ./...'
end

desc 'Run all tests'
task :test do
  sh 'go test ./...'
  sh 'go test -race ./...'
end

desc 'Run tests with coverage'
task :coverage do
  rm_f 'coverage.txt'
  coverage = ['mode: atomic']
  PACKAGES.each do |pkg|
    coverage += cover_one_package(pkg)
  end

  open('coverage.txt', 'w') { |f| coverage.each { |l| f.puts l } }
end

def cover_one_package(pkg)
  Dir.mktmpdir do |tdir|
    tcov = File.join(tdir, 'profile.out')

    sh(*%W(go test -coverprofile=#{tcov} -covermode=atomic #{pkg}))

    return [] unless File.exist? tcov

    new_lines = File.readlines(tcov)[2..-1]
    return [] if new_lines.nil?

    new_lines.compact.reject { |l| l =~ /\A\s*\Z/ }
  end
end

desc 'Builds the executables'
task :build do
  sh 'go build -v ./...'
end

desc 'Cross builds the executables'
task :xbuild do
  sh 'gox -os="linux darwin" -arch="amd64 ppc64le" ./...'
end

desc 'Cleans up executables'
task :clobber do
  rm FileList['docker-gc_*']
end
