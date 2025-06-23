async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("Deploying VIPPass with account:", deployer.address);

    const VIP = await ethers.getContractFactory("VIPPass");
    const vip = await VIP.deploy();
    await vip.deployed();

    console.log("VIPPass deployed to:", vip.address);
}
  
main().catch((err) => {
    console.error(err);
    process.exit(1);
});
  