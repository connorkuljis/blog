# copy local assets to remote staging directory
sudo rsync --delete -avz public/ $SSH_USER@$SERVER_IP:~/public/

# replace document root with staged directory assets
ssh -t $SSH_USER@$SERVER_IP \
	"sudo rm -rf /var/www/html/kuljis.xyz/* && sudo cp -r ~/public/* /var/www/html/kuljis.xyz/"

