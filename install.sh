sudo -s
wget https://go.dev/dl/go1.17.4.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.4.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
git clone https://github.com/cshi0/cs655-project.git
cd cs655-project
go build
