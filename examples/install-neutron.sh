#!/bin/bash

function setupNebulaDir(){
    sudo chown -R $(whoami) /etc/nebula
}

function download(){
    # Could handle multiple platforms here
    cd /tmp
    wget https://github.com/b177y/starship/releases/download/v0.3.0/$(uname -s)-$(uname -m).tar.gz -O /tmp/starship.tar.gz
    rm -r /tmp/release
    tar -xf /tmp/starship.tar.gz
    release=$(echo "/tmp/release/$(uname -s)-$(uname -m)" | tr '[:upper:]' '[:lower:]')
    sudo mv "${release}/nebula" /usr/local/bin/nebula
    sudo mv "${release}/neutron" /usr/local/bin/neutron
    sudo chown root:root /usr/local/bin/nebula
    sudo chown $(whoami):$(whoami) /usr/local/bin/neutron
    sudo mv "${release}/nebula@.service" /etc/systemd/system/nebula@.service
    sudo chown root:root /etc/systemd/system/nebula@.service
    sudo systemctl daemon-reload
    cd -
}

function runNeutron(){
    echo -n "Quasar Endpoint: "
    read qaddr
    echo -n "Network to Join: "
    read netname
    echo -n "Node Name: "
    read nodename
    /usr/local/bin/neutron join -quasar $qaddr -network $netname -name $nodename || exit
    echo "Run \`neutron update -network $netname\` to get the node config once you have approved the node."
    echo "Once the config has been fetched you can start nebula with \`systemctl start nebula@$netname\`"
}

function main(){
    download
    setupNebulaDir
    runNeutron
}

main
