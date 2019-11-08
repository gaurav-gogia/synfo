import sys
import glob
from fastai.vision.image import open_image
from fastai.basic_train import load_learner
from fastai.vision.transform import get_transforms
from fastai.vision.data import ImageDataBunch, imagenet_stats


def prep_model(path: str):
    model_path = './model/'
    classes = ('knife', 'handgun', 'rifle')

    data = ImageDataBunch.single_from_classes(
        path, classes, ds_tfms=get_transforms(), size=224).normalize(imagenet_stats)
    return load_learner(model_path)


def weapon_check(path_list: list, model):
    for path in path_list:
        detect(path, model)


def detect(img_path: str, model):
    img = open_image(img_path)
    pred_class, pred_idx, outputs = model.predict(img)

    if (outputs[0] < 0.7) and (outputs[1] < 0.7) and (outputs[2] < 0.7):
        print('#######################################################################')
        print('File: ', img_path.split('/')[-1])
        print('kai na malyu')
        print()
    else:
        print('#######################################################################')
        print('File: ' + img_path.split('/')[-1])
        print('prediction: ' + pred_class)
        print('probabilities:')
        print('handgun: ' + str(outputs[0]) + '\nknife: ' +
              str(outputs[1]) + '\nrifle: ' + str(outputs[2]))


def main():
    outfile = open('./weapon_rec_result.txt', 'w')
    path = sys.argv[1]

    if path == "":
        print("ERROR: Path cannot be empty", file=outfile)
        return

    model = prep_model(path)
    path_list = glob.glob(path+"*")
    weapon_check(path_list, model)


if __name__ == "__main__":
    main()
