// import {QtumProvider} from "./ethers/index.js";
// import * as a from "./ethers/index.js";
// import * as a from "ethers";
// import * as $ from "jQuery";
// import {$,jQuery} from 'jquery';
// import TruffleContract from "truffle-contract";
import 'regenerator-runtime/runtime'
import {providers, Contract, ethers} from "ethers"
import {QtumProvider, QtumWeb3Provider, QtumWallet} from "./ethers"
import {utils} from "web3"
console.log("QtumProvider", QtumProvider)
console.log("QtumWallet", QtumWallet)
var $ = require( "jquery" );
import AdoptionArtifact from './Adoption.json'
import Pets from './pets.json'
// export for others scripts to use
window.$ = $;
window.jQuery = $;

let QTUMMainnet = {
  chainId: '0x71',
  chainName: 'Qtum Mainnet',
  rpcUrls: ['https://localhost:23888'],
  blockExplorerUrls: ['https://qtum.info/'],
  iconUrls: [
    'https://qtum.info/images/metamask_icon.svg',
    'https://qtum.info/images/metamask_icon.png',
  ],
};
let QTUMTestNet = {
  chainId: '0x71',
  chainName: 'Qtum Testnet',
  rpcUrls: ['https://localhost:23888'],
  blockExplorerUrls: ['https://testnet.qtum.info/'],
  iconUrls: [
    'https://qtum.info/images/metamask_icon.svg',
    'https://qtum.info/images/metamask_icon.png',
  ],
};
let config = {
  "0x1": QTUMMainnet,
  // ETH Ropsten
  "0x3": QTUMTestNet,
  // ETH Rinkby
  "0x4": QTUMTestNet,
  // ETH Görli
  "0x5": QTUMTestNet,
  // ETH Kovan
  "0x71": QTUMTestNet,
};
config[QTUMMainnet.chainId] = QTUMMainnet;
config[QTUMTestNet.chainId] = QTUMTestNet;

window.App = {
  web3Provider: null,
  contracts: {},
  account: "",

  init: function() {
    // Load pets.
    var petsRow = $('#petsRow');
    var petTemplate = $('#petTemplate');

    for (let i = 0; i < Pets.length; i ++) {
      petTemplate.find('.panel-title').text(Pets[i].name);
      petTemplate.find('img').attr('src', Pets[i].picture);
      petTemplate.find('.pet-breed').text(Pets[i].breed);
      petTemplate.find('.pet-age').text(Pets[i].age);
      petTemplate.find('.pet-location').text(Pets[i].location);
      petTemplate.find('.btn-adopt').attr('pets-id', Pets[i].id);

      petsRow.append(petTemplate.html());
    }

    App.login()
    return App.initWeb3();
  },

  getChainId: function() {
    return window.ethereum.chainId;
  },
  isOnQtumChainId: function() {
    let chainId = this.getChainId();
    return chainId == QTUMMainnet.chainId || chainId == QTUMTestNet.chainId;
  },

  initWeb3: function() {
    let self = this;
    let qtumConfig = config[this.getChainId()] || QTUMTestNet;
    console.log("Adding network to Metamask", qtumConfig);
    window.ethereum.request({
      method: "wallet_addEthereumChain",
      params: [qtumConfig],
    })
      .then(() => {
        console.log("Successfully connected to qtum")
        window.ethereum.request({ method: 'eth_requestAccounts' })
          .then((accounts) => {
            console.log("Successfully logged into metamask", accounts);
            let qtumConnected = self.isOnQtumChainId();
            let currentlyQtumConnected = self.qtumConnected;
            if (accounts && accounts.length > 0) {
              App.account = accounts[0];
            }
            // self.$emit("web3Connected", true);
            if (currentlyQtumConnected != qtumConnected) {
              console.log("ChainID matches QTUM, not prompting to add network to web3, already connected.");
              // self.$emit("qtumConnected", true);
            }
            // let provider = new providers.Web3Provider(window.ethereum);
            let provider = new QtumWeb3Provider(window.ethereum);
            App.web3Provider = provider.getSigner()
            // const signer = provider.getSigner();
            console.log("provider", provider)
            // console.log("signer", signer);
            return App.initContract();
          })
          .catch((e) => {
            console.log("Connecting to web3 failed", e);
          })
      })
      .catch(() => {
        console.log("Adding network failed", arguments);
      })
    // App.web3Provider = new Web3.providers.HttpProvider('http://localhost:23889');
  },

  initContract: async function() {
    // App.contracts.Adoption = TruffleContract(AdoptionArtifact);
    let chainId = utils.hexToNumber(this.getChainId())
    console.log(chainId)
    App.contracts.Adoption = new Contract(AdoptionArtifact.networks[''+chainId].address, AdoptionArtifact.abi, App.web3Provider)

    // Set the provider for our contract
    // App.contracts.Adoption.setProvider(App.web3Provider);

    // Use our contract to retrieve and mark the adopted pets
    await App.markAdopted();
    return App.bindEvents();
  },

  bindEvents: function() {
    console.log("bindEvents")
    $(document).on('click', '.btn-adopt', App.handleAdopt);
  },

  markAdopted: function(adopters, account) {
    console.log("markAdopted")
    var adoptionInstance;
    return new Promise((resolve, reject) => {
      console.log("markAdopted promise")
      let deployed = App.contracts.Adoption.deployed();
      deployed.then(function(instance) {
        adoptionInstance = instance;
        console.log("instance", instance)
        return adoptionInstance.getAdopters.call()
          .then(function(adopters) {
            console.log("adopters", adopters)
            for (var i = 0; i < adopters.length; i++) {
              const adopter = adopters[i];
              if (adopter !== '0x0000000000000000000000000000000000000000') {
                $('.panel-pet').eq(i).find('button').text('Adopted').attr('disabled', true);
                $('.panel-pet').eq(i).find('.pet-adopter-container').css('display', 'block');
                let adopterLabel = adopter;
                if (adopter === App.account) {
                  adopterLabel = "You"
                }
                $('.panel-pet').eq(i).find('.pet-adopter-address').text(adopterLabel);
              } else {
                $('.panel-pet').eq(i).find('.pet-adopter-container').css('display', 'none');
              }
            }
            resolve()
            console.log("markAdopted success!")
          }).catch(function(err) {
            console.log(err);
            reject(err)
          });
      }).catch(function(err) {
        console.error(err)
      })
    });
  },

  handleAdopt: function(event) {
    event.preventDefault();

    var petId = parseInt($(event.target).data('id'));

    var adoptionInstance;

    console.log("handleAdopt")
    App.contracts.Adoption.deployed().then(function(instance) {
      adoptionInstance = instance;

      console.log("handleAdopt, adopt")
      return adoptionInstance.adopt(petId, {from: App.account});
    }).then(function(result) {
      console.log("handleAdopt success")
      return App.markAdopted();
    }).catch(function(err) {
      console.log("handleAdopt", err)
      console.log(err.message);
    });
  },

  login: function() {
    return
    let walletAddress = localStorage.getItem("userWalletAddress");
    while (!walletAddress) {
      walletAddress = window.prompt("Please enter your wallet address");
      if (walletAddress) {
        localStorage.setItem("userWalletAddress", walletAddress);
      }
    }

    App.account = walletAddress;
  },

  handleLogout: function() {
    localStorage.removeItem("userWalletAddress");

    App.login();
    App.markAdopted();
  }
};

$(function() {
  $(document).ready(function() {
    App.init();
  });
});
