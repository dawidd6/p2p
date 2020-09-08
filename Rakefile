binaries = %w[
  p2p
  p2pd
  p2p-tracker
  p2p-trackerd
]

protos = %w[
  client
  tracker
  metadata
]

task proto: :proto_plugins do
  protos.each do |proto|
    sh "protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/#{proto}/#{proto}.proto"
    sh "protoc --go_out=. --go_opt=paths=source_relative pkg/#{proto}/#{proto}.proto"
  end
end

task :proto_plugins do
  ENV["GO111MODULE"] = "off"
  sh "go get google.golang.org/protobuf/cmd/protoc-gen-go"
  sh "go get google.golang.org/grpc/cmd/protoc-gen-go-grpc"
end

task :default do
  binaries.each do |binary|
    sh "go build -o ./bin/#{binary} ./cmd/#{binary}"
  end
end