import sys
import time
import glob
import datetime
import face_recognition
from colorama import Fore, Style


def match(test, train_file_names):
    for train in train_file_names:
        print(Style.RESET_ALL)
        test_img = face_recognition.load_image_file(test)
        test_enc = face_recognition.face_encodings(test_img)

        train_img = face_recognition.load_image_file(train)
        train_enc = face_recognition.face_encodings(train_img)

        if (len(test_enc) > 0) and (len(train_enc) > 0):

            test_enc = test_enc[0]
            train_enc = train_enc[0]

            ans = face_recognition.compare_faces([test_enc], train_enc)
            if ans[0]:
                print(Fore.BLUE + test + ' = ' + train)
            else:
                print(Fore.RED + test + ' does NOT match ' + train)
        else:
            if len(test_enc) <= 0:
                print('No faces found in: ', test)
            if len(train_enc) <= 0:
                print('No faces found in: ', train)


def compare_images(test_file_names: list, train_file_names: list):
    x = 1
    for test in test_file_names:
        print(Style.RESET_ALL)
        print('#######################################################################')
        print('Phase: ', x)
        match(test, train_file_names)
        x += 1
        print('#######################################################################')


def main():
    start = time.time()

    test_dir = sys.argv[1]
    train_dir = sys.argv[2]

    if test_dir == "" or train_dir == "":
        print("Paths cannot be empty")
    else:
        test_file_names = glob.glob(test_dir+"*")
        train_file_names = glob.glob(train_dir+"*")

        compare_images(test_file_names, train_file_names)

        end = time.time()
        print(Style.RESET_ALL)
        print('\nElapsed Time: ', str(datetime.timedelta(seconds=end - start)))


if __name__ == "__main__":
    main()
