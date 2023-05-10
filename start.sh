make down
sleep 5
docker system prune --force
docker ps
make up
sleep 5
make private_key=1ec1d7b96ee2512390f755475379448f1a0886cc3b050ed09da58e95329b2fd0 add-private-key
sleep 5
make add-vault-token
sleep 5
make generate-issuer-did
sleep 5
make run
sleep 5
make run-ui
sleep 5
docker ps
