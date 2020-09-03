import json


def main():
    with open('result/gctest/nogc.json') as f:
        data = json.load(f)
    nogc = [item[0] for item in data]
    with open('result/gctest/gc1min.json') as f:
        data = json.load(f)
    gc1 = [item[0] for item in data]
    with open('result/gctest/gc10min.json') as f:
        data = json.load(f)
    gc10 = [item[0] for item in data]
    with open('result/gctest/txn.json') as f:
        data = json.load(f)
    txn = [item[0] for item in data]
    with open('result/gctest/gc', 'w') as f:
        f.write("#{:<19} {:<20} {:<20} {:<20} {:<20}\n".format("timepoint", "Without GC", "GC(T=1 min)", "GC(10 min)",
                                                             "Beldi-Txn"))
        for i in range(31):
            try:
                f.write("{:<20} {:<20} {:<20} {:<20} {:<20}\n".format(i, nogc[i], gc1[i], gc10[i], txn[i]))
            except:
                break


if __name__ == '__main__':
    main()
