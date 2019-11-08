import face_recognition
import datetime
import glob
import time
import sys


def match(test: str, train_file_names: list, outfile, model_type: str):
    for train in train_file_names:
        test_img = face_recognition.load_image_file(test)
        test_loc = face_recognition.face_locations(test_img, model=model_type)
        test_enc = face_recognition.face_encodings(
            test_img, known_face_locations=test_loc, num_jitters=5, model_size="large")

        train_img = face_recognition.load_image_file(train)
        train_loc = face_recognition.face_locations(
            train_img, model=model_type)
        train_enc = face_recognition.face_encodings(
            train_img, known_face_locations=train_loc, num_jitters=5, model_size="large")

        if (len(test_enc) == 1) and (len(train_enc) == 1):
            test_enc = test_enc[0]
            train_enc = train_enc[0]
            ans = face_recognition.compare_faces([test_enc], train_enc)
            if ans[0]:
                print(test + ' = ' + train, file=outfile)
            else:
                print(test + ' does NOT match ' + train, file=outfile)

        elif(len(test_enc) > 1 and len(train_enc) > 1):
            match_both_phase1(test, train, outfile, test_enc, train_enc)

        elif (len(test_enc) > 1 and len(train_enc) == 1):
            train_enc = train_enc[0]
            match_one_multi(test, train, outfile, train_enc, test_enc, False)

        elif(len(train_enc) > 1 and len(test_enc) == 1):
            test_enc = test_enc[0]
            match_one_multi(test, train, outfile, test_enc, train_enc, True)

        else:
            if len(test_enc) <= 0:
                print('No faces found in: ' + test, file=outfile)
            if len(train_enc) <= 0:
                print('No faces found in: ' + train, file=outfile)


def match_one_multi(test: str, train: str, outfile: str, encsingle, enclist, train_multi: bool):
    i = 1
    for enc in enclist:
        if train_multi:
            fname = train + '(face ' + str(i) + ') from left'
        else:
            fname = test + '(face ' + str(i) + ') from left'
        ans = face_recognition.compare_faces([enc], encsingle)
        if ans[0]:
            if train_multi:
                print(test + ' from left = ' + fname, file=outfile)
            else:
                print(fname + ' from left = ' + train, file=outfile)

        else:
            if train_multi:
                print(test + ' does NOT match ' +
                      fname, file=outfile)
            else:
                print(fname + ' does NOT match ' +
                      train, file=outfile)
        i += 1


def match_both_phase1(test: str, train: str, outfile: str, test_list, train_list):
    i = 1
    for tenc in test_list:
        match_both_phase2(test, train, outfile, i, tenc, train_list)
        i += 1


def match_both_phase2(test: str, train: str, outfile: str, test_findex: int, test_enc, train_list):
    i = 1
    testname = test + '(face ' + str(test_findex) + ')'
    for trinc in train_list:
        trainame = train + '(face ' + str(i) + ')'
        ans = face_recognition.compare_faces([test_enc], trinc)

        if ans[0]:
            print(testname + ' from left = ' +
                  trainame + ' from left ', outfile)
        else:
            print(testname + ' from left does NOT match ' +
                  trainame + ' from left ', outfile)
        i += 1


def compare_images(test_file_names: list, train_file_names: list, outfile, model_type: str):
    x = 1
    for test in test_file_names:
        print('#######################################################################', file=outfile)
        print('Phase: ', x, file=outfile)

        match(test, train_file_names, outfile, model_type)
        x += 1
        print('#######################################################################', file=outfile)


def main():
    outfile = open('./face_rec_result.txt', 'w')

    test_dir = sys.argv[1]
    train_dir = sys.argv[2]
    model_type = sys.argv[3]

    if test_dir == "" or train_dir == "":
        print("ERROR: Paths cannot be empty", file=outfile)
    elif model_type == "cnn" or model_type == "hog":
        test_file_names = glob.glob(test_dir+"*")
        train_file_names = glob.glob(train_dir+"*")

        compare_images(test_file_names, train_file_names, outfile, model_type)
    else:
        print("ERROR: Model type can only be either cnn or hog", file=outfile)

    outfile.close()


if __name__ == "__main__":
    main()
