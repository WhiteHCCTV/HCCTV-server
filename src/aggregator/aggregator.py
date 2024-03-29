from datetime import datetime
from torch import nn, optim
import torch
from models import models
from utils.config import Config
from utils.data import get_data
from torch.utils.data import DataLoader
import numpy as np
import logging
import pickle
import os
import copy

import io
class CPU_Unpickler(pickle.Unpickler):
    def find_class(self, module, name):
        if module == 'torch.storage' and name == '_load_from_bytes':
            return lambda b: torch.load(io.BytesIO(b), map_location='cpu')
        else:
            return super().find_class(module, name)

class Updater:
    def __init__(self, config):
        self.config = config
        self.global_weights = []
        self.model = self.load_model()
        self.test_loader = self.load_test_data()
        self.global_weights = None
        self.init_weights = None

    def load_model(self):
        # logging.info('Load Model: {}'.format(self.config.dataset))
        print(datetime.now().strftime("[%H:%M:%S]")+'Load Model: {}'.format(self.config.dataset))
        
        model = models.get_model()
        return model

    def load_test_data(self):
        _, testset = get_data(self.config.dataset, self.config)
        test_loader = DataLoader(testset, batch_size=32, shuffle=True)
        return test_loader

    def test(self, is_init):
        # set_weights
        if(is_init):
            self.model.load_state_dict(self.init_weights)
        else:
            self.model.load_state_dict(self.global_weights)
        
        corrects = 0
        test_loss = 0
        device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
        self.model.to(device)
        self.model.eval()
        criterion = nn.CrossEntropyLoss()

        #dataloader = DataLoader(self.testset, batch_size=32, shuffle=True)
        for batch_id, (inputs, labels) in enumerate(self.test_loader):
            loss = 0
            inputs = inputs.to(device)
            labels = labels.to(device)
            outputs = self.model(inputs)
            loss = criterion(outputs, labels)
            _, preds = torch.max(outputs, 1)
            test_loss += loss.item() * inputs.size(0)
            corrects += torch.sum(preds == labels.data)

        acc = int(corrects) / len(self.test_loader.dataset)
        avg_loss = test_loss / len(self.test_loader.dataset)

        # logging.info("Test Accuracy: {}, Avgerage Loss: {}".format(acc, avg_loss))
        print(datetime.now().strftime("[%H:%M:%S]")+"Test Accuracy: {}, Avgerage Loss: {}".format(acc, avg_loss))

        return acc, avg_loss

    def fed_avg(self, dir_path):
        # get_file_list of Dir
        local_weights = []
        file_list = os.listdir(dir_path)

        client_num = len(file_list)
        # logging.info("Client_Num: {}".format(client_num))
        print(datetime.now().strftime("[%H:%M:%S]")+"Client_Num: {}".format(client_num))

        for i in range(client_num):
            file_path = dir_path + file_list[i]
            local_weights.append(self.load_weights(file_path))
        
        w_avg = copy.deepcopy(local_weights[0])
        for k in w_avg.keys():
            for i in range(1, client_num):
                w_avg[k] = torch.add(w_avg[k], local_weights[i][k])
            w_avg[k] = torch.div(w_avg[k], float(client_num))
        self.global_weights = w_avg

        # Save aggregated global weights
        self.save_weights(self.global_weights, './aggregator/globals/global.pickle')

    def save_weights(self, givenWeights, PATH):
        with open(PATH, 'wb') as f:
            pickle.dump(givenWeights, f)

    def load_weights(self, filePath):
        # Load local weights from .pickle
        with open(filePath, 'rb') as inputfile:
            #weights = pickle.load(inputfile)
            weights = CPU_Unpickler(inputfile).load()
        return weights

    def set_init_weights(self, filePath):
        with open(filePath, 'rb') as inputfile:
        weights = CPU_Unpickler(inputfile).load()
        # logging.info("Init Weights Test")
        print(datetime.now().strftime("[%H:%M:%S]")+"Init Weights Test")        
        self.init_weights = weights

if __name__=="__main__":
    logging.basicConfig(
        format='[%(levelname)s][%(asctime)s]: %(message)s',
        level=getattr(logging, "INFO"), datefmt='%H:%M:%S'
    )
    PATH = './aggregator/locals/'
    lst = os.listdir(PATH)
    config = Config("./aggregator/configs/params.json")

    initTester = Updater(config)
    initTester.set_init_weights(PATH+str(lst[0]))
    # logging.info("Init Weights Accuracy using '{}'".format(lst[0]))
    print(datetime.now().strftime("[%H:%M:%S]")+"Init Weights Accuracy using '{}'".format(lst[0]))
    initTester.test(True)

    print()

    updater = Updater(config)
    # logging.info("Global Weights Accuracy")
    print(datetime.now().strftime("[%H:%M:%S]")+"Global Weights Accuracy")
    updater.fed_avg(PATH)
    updater.test(False)
