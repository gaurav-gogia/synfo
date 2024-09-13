import os
import sys
import glob
from fastai.vision.all import *

def load_model(path: str):
    cwd = os.getcwd()
    model_path = cwd + "/" + "model/awd.pkl"
    return load_learner(model_path)


def weapon_check(path_list: list, model, out):
    for path in path_list:
        print(path)
        detect(path, model, out)


def detect(img_path: str, model, out):
    img = load_image(img_path)
    pred_class, pred_idx, outputs = model.predict(img)

    if (outputs[0] < 0.7) and (outputs[1] < 0.7) and (outputs[2] < 0.7):
        print("#######################################################################", file=out)
        print("File: ", img_path.split("/")[-1], file=out)
        print("no weapon found", file=out)
        print("", file=out)
    else:
        print("#######################################################################", file=out)
        print("File: " + img_path.split("/")[-1], file=out)
        print("prediction: " + pred_class, file=out)
        print("probabilities:", file=out)
        print("handgun: " + str(outputs[0]) + "\nknife: " + str(outputs[1]) + "\nrifle: " + str(outputs[2]), file=out)


def main():
    cwd = os.getcwd()
    outpath = cwd + "/weapon_rec_result.txt"
    outfile = open(outpath, "w")

    if len(sys.argv) < 2:
        print("ERROR: Path cannot be empty", file=outfile)
        return

    path = sys.argv[1]
    model = load_model(path)
    path_list = glob.glob(path+"*")
    weapon_check(path_list, model, outfile)


if __name__ == "__main__":
    main()
