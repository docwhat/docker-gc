task default: %w(lint test build)

desc 'Fetches development and build dependencies'
task :setup do
  sh 'go get -u github.com/alecthomas/gometalinter'
  sh 'gometalinter --install --update'
  sh 'go get -u github.com/mitchellh/gox'
end

desc 'Check lint and style'
task :lint do
  sh 'gometalinter --disable=gotype ./...'
end

desc 'Run all tests'
task :test do
  sh 'go test -v ./...'
  sh 'go test -v -race ./...'
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
