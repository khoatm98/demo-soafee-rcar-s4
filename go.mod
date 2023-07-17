module github.com/khoatm98/demo-soafee-rcar-s4

go 1.18

replace (
	github.com/ThalesIgnite/crypto11 => github.com/aoscloud/crypto11 v1.0.3-0.20220217163524-ddd0ace39e6f
	github.com/syucream/posix_mq => github.com/al1img/posix_mq v0.0.2-0.20220603145914-6cbbc81f1d84
)

require (
	github.com/aoscloud/aos_common v0.0.0-20220818090503-b3b09ab17df8
	github.com/aoscloud/aos_updatemanager v0.0.0-20220818090328-00a6b97ce9a2
	github.com/sirupsen/logrus v1.8.1
	github.com/syucream/posix_mq v0.0.1
)

require (
	github.com/ThalesIgnite/crypto11 v0.0.0-00010101000000-000000000000 // indirect
	github.com/cavaliergopher/grab/v3 v3.0.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-tpm v0.3.3 // indirect
	github.com/looplab/fsm v0.3.0 // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/thales-e-security/pool v0.0.2 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220317061510-51cd9980dadf // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106 // indirect
	google.golang.org/grpc v1.46.2 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
