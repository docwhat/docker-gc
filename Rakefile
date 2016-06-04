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
  sh 'gometalinter --deadline=20s --disable=gotype ./...'
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
  PACKAGES.each do |path|
    sh(*%W(go test -coverprofile=profile.out -covermode=atomic #{path}))
    if File.exist? 'profile.out'
      coverage << File.readlines('profile.out')[2..-1]
      rm_f 'profile.out'
    end
  end
  open('coverage.txt', 'w') do |f|
    coverage.flatten.each { |l| f.puts l }
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
