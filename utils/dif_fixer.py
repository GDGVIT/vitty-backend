from difflib import SequenceMatcher as sq


def fix_string(arg):
    arg = arg.replace("g","0")
    arg = arg.replace("c","0")
    arg = arg.replace("a", "3")
    arg = arg.replace("o", "0")
    arg = arg.replace('w',"V")
    arg = arg.replace("h","1")
    arg.upper()
    if len(arg) <= 4:
        if "I" in arg:
            arg = arg.replace("I", "1")
        elif arg.startswith("1"):
            arg = arg.replace("1","L",1)
        elif arg.endswith("T"):
            arg = arg.replace(arg,arg[:-1]+"1",1)
        elif arg.endswith("D"):
            arg = arg.replace(arg,arg[:-1]+"0",1)

    if len(arg) >= 4:
        if "O" in arg:
            arg = arg.replace("O", "0")
        elif "B" in arg[3:]:
            arg = arg.replace("B", "8")
        elif "I" in arg[3:]:
            arg = arg.replace("I", "1")
        elif "E" in arg[3:]:
            arg = arg.replace("E", "8")
        elif "K" in arg[2:]:
            arg = arg.replace("K", "1")
        elif arg.startswith("5"):
            arg = arg.replace("5","S",1)
    return arg
