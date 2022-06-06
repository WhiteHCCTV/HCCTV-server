# TO DO #
from torchvision import datasets, transforms
from torch.utils.data import DataLoader
import numpy as np

class MNIST:
    def __init__(self, config):
        self.config = config
        print(self.config.paths.data)
        self.path = str(self.config.paths.data) + "/" + self.config.dataset
        self.trainset = None
        self.testset = None

    def load_data(self):
        self.trainset = datasets.MNIST(
            self.path, train=True, download=True, transform=transforms.Compose([
                transforms.RandomRotation(15),
                transforms.ToTensor(),
                transforms.Normalize((0.1307,), (0.3081,))
            ])
        )
        self.testset = datasets.MNIST(
            self.path, train=False, download=True, transform=transforms.Compose([
                transforms.RandomRotation(15),
                transforms.ToTensor(),
                transforms.Normalize((0.1307,), (0.3081,))
            ])
        )
        return self.trainset, self.testset

class FashionMNIST:
    def __init__(self, config):
        self.config = config
        print(self.config.paths.data)
        self.path = str(self.config.paths.data) + "/" + self.config.dataset
        self.trainset = None
        self.testset = None

    def load_data(self):
        self.trainset = datasets.FashionMNIST(
            self.path, train=True, download=True, transform=transforms.Compose([
                transforms.RandomRotation(15),
                transforms.ToTensor(),
                transforms.Normalize((0.1307,), (0.3081,))
            ])
        )
        self.testset = datasets.FashionMNIST(
            self.path, train=False, download=True, transform=transforms.Compose([
                transforms.RandomRotation(15),
                transforms.ToTensor(),
                transforms.Normalize((0.1307,), (0.3081,))
            ])
        )
        return self.trainset, self.testset

def get_data(dataset, config):
    if dataset == "MNIST":
        return MNIST(config).load_data()
    elif dataset == "FashionMNIST":
        return FashionMNIST(config).load_data()

