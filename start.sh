make down
sleep 10
docker system prune --force
docker ps
make up
sleep 10
make private_key=1ec1d7b96ee2512390f755475379448f1a0886cc3b050ed09da58e95329b2fd0 add-private-key
sleep 10
make add-vault-token
sleep 10
make generate-issuer-did
sleep 10
make run
sleep 10
make run-ui
sleep 10
docker ps
