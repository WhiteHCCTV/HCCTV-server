import torch.nn.functional as F
from torch import nn

def get_model():
    return TORCH_NET()

class TORCH_NET(nn.Module):
    def __init__(self):
        super(TORCH_NET, self).__init__()
        self.conv1 = nn.Conv2d(1, 32, 5)
        self.pool = nn.MaxPool2d(2,2)
        self.conv2 = nn.Conv2d(32, 64, 5)
        self.fc = nn.Sequential(
                nn.Linear(64 * 4 * 4, 512),
                nn.ReLU(),
                nn.Linear(512, 10)
        )

    def forward(self, x):
        x = self.pool(F.relu(self.conv1(x)))
        x = self.pool(F.relu(self.conv2(x)))
        x = x.view(-1, 4*4*64)
        x = self.fc(x)
        return x

