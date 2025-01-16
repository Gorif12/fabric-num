#!/bin/bash

# 添加错误处理函数
handle_error() {
    local exit_code=$?
    local step_name=$1
    echo "❌ 步骤失败: $step_name"
    echo "错误代码: $exit_code"
    exit $exit_code
}

# 添加进度显示函数
show_progress() {
    local current_step=$1
    local total_steps=15
    local step_name=$2
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🔄 [步骤 $current_step/$total_steps] $step_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

ChannelName="channeltogettoyou"
ChainCodeName="chaincodetogettoyou"

# base
BASE_PATH="/etc/hyperledger/peer"
ORDERER1_ADDRESS="orderer1.togettoyou.com:7050"
ORDERER_CA="/etc/hyperledger/orderer/togettoyou.com/orderers/orderer1.togettoyou.com/msp/tlscacerts/tlsca.togettoyou.com-cert.pem"
CORE_PEER_TLS_ENABLED=true

# Org1 Peer0
ORG1_PEER0_ADDRESS="peer0.org1.togettoyou.com:7051"
ORG1_PEER0_LOCALMSPID="Org1MSP"
ORG1_PEER0_MSPCONFIGPATH="${BASE_PATH}/org1.togettoyou.com/users/Admin@org1.togettoyou.com/msp"
ORG1_PEER0_TLS_ROOTCERT_FILE="${BASE_PATH}/org1.togettoyou.com/peers/peer0.org1.togettoyou.com/tls/ca.crt"
ORG1_PEER0_TLS_CERT_FILE="${BASE_PATH}/org1.togettoyou.com/peers/peer0.org1.togettoyou.com/tls/server.crt"
ORG1_PEER0_TLS_KEY_FILE="${BASE_PATH}/org1.togettoyou.com/peers/peer0.org1.togettoyou.com/tls/server.key"

# Org1 Peer1
ORG1_PEER1_ADDRESS="peer1.org1.togettoyou.com:7051"
ORG1_PEER1_LOCALMSPID="Org1MSP"
ORG1_PEER1_MSPCONFIGPATH="${BASE_PATH}/org1.togettoyou.com/users/Admin@org1.togettoyou.com/msp"
ORG1_PEER1_TLS_ROOTCERT_FILE="${BASE_PATH}/org1.togettoyou.com/peers/peer1.org1.togettoyou.com/tls/ca.crt"
ORG1_PEER1_TLS_CERT_FILE="${BASE_PATH}/org1.togettoyou.com/peers/peer1.org1.togettoyou.com/tls/server.crt"
ORG1_PEER1_TLS_KEY_FILE="${BASE_PATH}/org1.togettoyou.com/peers/peer1.org1.togettoyou.com/tls/server.key"

# Org2 Peer0
ORG2_PEER0_ADDRESS="peer0.org2.togettoyou.com:7051"
ORG2_PEER0_LOCALMSPID="Org2MSP"
ORG2_PEER0_MSPCONFIGPATH="${BASE_PATH}/org2.togettoyou.com/users/Admin@org2.togettoyou.com/msp"
ORG2_PEER0_TLS_ROOTCERT_FILE="${BASE_PATH}/org2.togettoyou.com/peers/peer0.org2.togettoyou.com/tls/ca.crt"
ORG2_PEER0_TLS_CERT_FILE="${BASE_PATH}/org2.togettoyou.com/peers/peer0.org2.togettoyou.com/tls/server.crt"
ORG2_PEER0_TLS_KEY_FILE="${BASE_PATH}/org2.togettoyou.com/peers/peer0.org2.togettoyou.com/tls/server.key"

# Org2 Peer1
ORG2_PEER1_ADDRESS="peer1.org2.togettoyou.com:7051"
ORG2_PEER1_LOCALMSPID="Org2MSP"
ORG2_PEER1_MSPCONFIGPATH="${BASE_PATH}/org2.togettoyou.com/users/Admin@org2.togettoyou.com/msp"
ORG2_PEER1_TLS_ROOTCERT_FILE="${BASE_PATH}/org2.togettoyou.com/peers/peer1.org2.togettoyou.com/tls/ca.crt"
ORG2_PEER1_TLS_CERT_FILE="${BASE_PATH}/org2.togettoyou.com/peers/peer1.org2.togettoyou.com/tls/server.crt"
ORG2_PEER1_TLS_KEY_FILE="${BASE_PATH}/org2.togettoyou.com/peers/peer1.org2.togettoyou.com/tls/server.key"

# Org1 Peer0 CLI 配置
Org1Peer0Cli="CORE_PEER_ADDRESS=${ORG1_PEER0_ADDRESS} \
CORE_PEER_LOCALMSPID=${ORG1_PEER0_LOCALMSPID} \
CORE_PEER_MSPCONFIGPATH=${ORG1_PEER0_MSPCONFIGPATH} \
CORE_PEER_TLS_ENABLED=${CORE_PEER_TLS_ENABLED} \
CORE_PEER_TLS_ROOTCERT_FILE=${ORG1_PEER0_TLS_ROOTCERT_FILE} \
CORE_PEER_TLS_CERT_FILE=${ORG1_PEER0_TLS_CERT_FILE} \
CORE_PEER_TLS_KEY_FILE=${ORG1_PEER0_TLS_KEY_FILE}"

# Org1 Peer1 CLI 配置
Org1Peer1Cli="CORE_PEER_ADDRESS=${ORG1_PEER1_ADDRESS} \
CORE_PEER_LOCALMSPID=${ORG1_PEER1_LOCALMSPID} \
CORE_PEER_MSPCONFIGPATH=${ORG1_PEER1_MSPCONFIGPATH} \
CORE_PEER_TLS_ENABLED=${CORE_PEER_TLS_ENABLED} \
CORE_PEER_TLS_ROOTCERT_FILE=${ORG1_PEER1_TLS_ROOTCERT_FILE} \
CORE_PEER_TLS_CERT_FILE=${ORG1_PEER1_TLS_CERT_FILE} \
CORE_PEER_TLS_KEY_FILE=${ORG1_PEER1_TLS_KEY_FILE}"

# Org2 Peer0 CLI 配置
Org2Peer0Cli="CORE_PEER_ADDRESS=${ORG2_PEER0_ADDRESS} \
CORE_PEER_LOCALMSPID=${ORG2_PEER0_LOCALMSPID} \
CORE_PEER_MSPCONFIGPATH=${ORG2_PEER0_MSPCONFIGPATH} \
CORE_PEER_TLS_ENABLED=${CORE_PEER_TLS_ENABLED} \
CORE_PEER_TLS_ROOTCERT_FILE=${ORG2_PEER0_TLS_ROOTCERT_FILE} \
CORE_PEER_TLS_CERT_FILE=${ORG2_PEER0_TLS_CERT_FILE} \
CORE_PEER_TLS_KEY_FILE=${ORG2_PEER0_TLS_KEY_FILE}"

# Org2 Peer1 CLI 配置
Org2Peer1Cli="CORE_PEER_ADDRESS=${ORG2_PEER1_ADDRESS} \
CORE_PEER_LOCALMSPID=${ORG2_PEER1_LOCALMSPID} \
CORE_PEER_MSPCONFIGPATH=${ORG2_PEER1_MSPCONFIGPATH} \
CORE_PEER_TLS_ENABLED=${CORE_PEER_TLS_ENABLED} \
CORE_PEER_TLS_ROOTCERT_FILE=${ORG2_PEER1_TLS_ROOTCERT_FILE} \
CORE_PEER_TLS_CERT_FILE=${ORG2_PEER1_TLS_CERT_FILE} \
CORE_PEER_TLS_KEY_FILE=${ORG2_PEER1_TLS_KEY_FILE}"

# 检查操作系统类型
if [[ `uname` == 'Darwin' ]]; then
  echo "当前操作系统是 Mac"
  export PATH=${PWD}/hyperledger-fabric-darwin-arm64-2.5.10/bin:$PATH
elif [[ `uname` == 'Linux' ]]; then
  echo "当前操作系统是 Linux"
  export PATH=${PWD}/hyperledger-fabric-linux-amd64-2.5.10/bin:$PATH
else
  echo "当前操作系统不是 Mac 或 Linux，脚本无法继续执行！"
  exit 1
fi

echo -e "注意：倘若您之前已经部署过了 network ，执行该脚本会丢失之前的数据！"
read -p "你确定要继续执行吗？请输入 Y 或 y 继续执行：" confirm

if [[ "$confirm" != "Y" && "$confirm" != "y" ]]; then
  echo "你取消了脚本的执行。"
  exit 1
fi

show_progress 1 "清理环境"
./uninstall.sh || handle_error "清理环境"

# 参考 https://hyperledger-fabric.readthedocs.io/zh-cn/release-2.5/test_network.html#id10

show_progress 2 "生成证书和秘钥（MSP 材料）"
cryptogen generate --config=./crypto-config.yaml || handle_error "生成证书和秘钥"

show_progress 3 "创建排序通道创世区块"
configtxgen -profile SampleGenesis -outputBlock ./config/genesis.block -channelID firstchannel || handle_error "创建排序通道创世区块"

show_progress 4 "生成通道配置事务"
configtxgen -profile SampleChannel -outputCreateChannelTx ./config/$ChannelName.tx -channelID $ChannelName || handle_error "生成通道配置事务"

show_progress 5 "定义组织锚节点"
configtxgen -profile SampleChannel -outputAnchorPeersUpdate ./config/Org1Anchor.tx -channelID $ChannelName -asOrg Org1 || handle_error "定义Org1锚节点"
configtxgen -profile SampleChannel -outputAnchorPeersUpdate ./config/Org2Anchor.tx -channelID $ChannelName -asOrg Org2 || handle_error "定义Org2锚节点"

show_progress 6 "启动所有节点"
docker-compose up -d || handle_error "启动节点"
echo "正在等待节点的启动完成，等待10秒"
sleep 10

CLI_CMD="docker exec cli.togettoyou.com bash -c"
CONFIG_PATH="/etc/hyperledger/config/"

show_progress 7 "创建通道"
$CLI_CMD "$Org1Peer0Cli peer channel create --outputBlock '$CONFIG_PATH'$ChannelName.block -o $ORDERER1_ADDRESS -c $ChannelName -f '$CONFIG_PATH'$ChannelName.tx --tls --cafile $ORDERER_CA" || handle_error "创建通道"

show_progress 8 "节点加入通道"
$CLI_CMD "$Org1Peer0Cli peer channel join -b '$CONFIG_PATH'$ChannelName.block" || handle_error "Org1Peer0加入通道"
$CLI_CMD "$Org1Peer1Cli peer channel join -b '$CONFIG_PATH'$ChannelName.block" || handle_error "Org1Peer1加入通道"
$CLI_CMD "$Org2Peer0Cli peer channel join -b '$CONFIG_PATH'$ChannelName.block" || handle_error "Org2Peer0加入通道"
$CLI_CMD "$Org2Peer1Cli peer channel join -b '$CONFIG_PATH'$ChannelName.block" || handle_error "Org2Peer1加入通道"

show_progress 9 "更新锚节点"
$CLI_CMD "$Org1Peer0Cli peer channel update -o $ORDERER1_ADDRESS -c $ChannelName -f '$CONFIG_PATH'Org1Anchor.tx --tls --cafile $ORDERER_CA" || handle_error "更新Org1锚节点"
$CLI_CMD "$Org2Peer0Cli peer channel update -o $ORDERER1_ADDRESS -c $ChannelName -f '$CONFIG_PATH'Org2Anchor.tx --tls --cafile $ORDERER_CA" || handle_error "更新Org2锚节点"

# 参考 https://hyperledger-fabric.readthedocs.io/zh-cn/release-2.5/deploy_chaincode.html

show_progress 10 "打包链码"
$CLI_CMD "peer lifecycle chaincode package /opt/gopath/src/chaincode/togettoyou_chaincode_1.0.0.tar.gz --path /opt/gopath/src/chaincode --lang golang --label togettoyou_chaincode_1.0.0" || handle_error "打包链码"

show_progress 11 "安装链码到对等节点"
$CLI_CMD "$Org1Peer0Cli peer lifecycle chaincode install /opt/gopath/src/chaincode/togettoyou_chaincode_1.0.0.tar.gz" || handle_error "Org1Peer0安装链码"
$CLI_CMD "$Org1Peer1Cli peer lifecycle chaincode install /opt/gopath/src/chaincode/togettoyou_chaincode_1.0.0.tar.gz" || handle_error "Org1Peer1安装链码"
$CLI_CMD "$Org2Peer0Cli peer lifecycle chaincode install /opt/gopath/src/chaincode/togettoyou_chaincode_1.0.0.tar.gz" || handle_error "Org2Peer0安装链码"
$CLI_CMD "$Org2Peer1Cli peer lifecycle chaincode install /opt/gopath/src/chaincode/togettoyou_chaincode_1.0.0.tar.gz" || handle_error "Org2Peer1安装链码"

show_progress 12 "批准链码定义"
Version="1.0.0"
Sequence="1"
PackageID=$($CLI_CMD "$Org1Peer0Cli peer lifecycle chaincode calculatepackageid /opt/gopath/src/chaincode/togettoyou_chaincode_1.0.0.tar.gz") || handle_error "计算PackageID"

$CLI_CMD "$Org1Peer0Cli peer lifecycle chaincode approveformyorg -o $ORDERER1_ADDRESS --channelID $ChannelName --name $ChainCodeName --version $Version --package-id $PackageID --sequence $Sequence --tls --cafile $ORDERER_CA" || handle_error "Org1批准链码"
$CLI_CMD "$Org1Peer0Cli peer lifecycle chaincode checkcommitreadiness --channelID $ChannelName --name $ChainCodeName --version $Version --sequence $Sequence --output json" || handle_error "Org1检查链码就绪状态"

$CLI_CMD "$Org2Peer0Cli peer lifecycle chaincode approveformyorg -o $ORDERER1_ADDRESS --channelID $ChannelName --name $ChainCodeName --version $Version --package-id $PackageID --sequence $Sequence --tls --cafile $ORDERER_CA" || handle_error "Org2批准链码"
$CLI_CMD "$Org2Peer0Cli peer lifecycle chaincode checkcommitreadiness --channelID $ChannelName --name $ChainCodeName --version $Version --sequence $Sequence --output json" || handle_error "Org2检查链码就绪状态"

show_progress 13 "提交链码定义"
$CLI_CMD "$Org1Peer0Cli peer lifecycle chaincode commit -o $ORDERER1_ADDRESS --channelID $ChannelName --name $ChainCodeName --version $Version --sequence $Sequence --tls --cafile $ORDERER_CA --peerAddresses $ORG1_PEER0_ADDRESS --tlsRootCertFiles $ORG1_PEER0_TLS_ROOTCERT_FILE --peerAddresses $ORG2_PEER0_ADDRESS --tlsRootCertFiles $ORG2_PEER0_TLS_ROOTCERT_FILE" || handle_error "提交链码定义"
$CLI_CMD "$Org1Peer0Cli peer lifecycle chaincode querycommitted --channelID $ChannelName --name $ChainCodeName" || handle_error "查询已提交的链码"

show_progress 14 "初始化链码"
$CLI_CMD "$Org1Peer0Cli peer chaincode invoke -o $ORDERER1_ADDRESS -C $ChannelName -n $ChainCodeName -c '{\"function\":\"InitLedger\",\"Args\":[]}' --tls --cafile $ORDERER_CA --peerAddresses $ORG1_PEER0_ADDRESS --tlsRootCertFiles $ORG1_PEER0_TLS_ROOTCERT_FILE --peerAddresses $ORG2_PEER0_ADDRESS --tlsRootCertFiles $ORG2_PEER0_TLS_ROOTCERT_FILE" || handle_error "初始化链码"

echo "正在等待链码初始化完成，等待5秒"
sleep 5

show_progress 15 "验证链码"
if $CLI_CMD "$Org1Peer0Cli peer chaincode query -C $ChannelName -n $ChainCodeName -c '{\"Args\":[\"Hello\"]}'" 2>&1 | grep "hello"; then
    echo "✅ 【恭喜您！】 network 部署成功"
    exit 0
fi

echo "❌ 【警告】network 未部署成功，请检查每一个步骤，定位具体问题。"
exit 1
