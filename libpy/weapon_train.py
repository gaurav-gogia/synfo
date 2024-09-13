import sys
import os
from fastai.vision.all import *

if len(sys.argv) < 2:
    print("Dir with training dataset required")
    return

print("Starting")

def train_model(path: str):
    classes = ['knife', 'handgun', 'rifle']
    dls = ImageDataLoaders.from_folder(
        path,
        valid_pct=0.2,
        item_tfms=Resize(224),
        batch_tfms=aug_transforms(),
        classes=classes
    )
    learn = vision_learner(dls, resnet34, metrics=accuracy)
    learn.fine_tune(4)
    cwd = os.getcwd()
    model_path = cwd + "/" + "model/awd.pkl"
    learn.export(model_path)

train_model(sys.argv[1])
print("Done!")