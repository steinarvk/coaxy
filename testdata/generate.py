import random
import json

def generate_sequence(n=1000):
    random.seed(1234)
    for _ in range(n):
        x = random.random()
        y = int(1000 * (x + 0.1 * (random.random()*2-1)))
        y *= y
        yield {
            "real": x,
            "integer": y,
        }

def generate_ctsv(f, separator, header, seq, keys):
    if header:
        f.write(separator.join(keys))
        f.write("\n")
    for rec in seq:
        f.write(separator.join([str(rec[k]) for k in keys]))
        f.write("\n")

def generate_logfmt(f, seq, keys):
    for rec in seq:
        f.write(" ".join([f"{k}={str(rec[k])}" for k in keys]))
        f.write("\n")

if __name__ == "__main__":
    name = "simpletest"
    keys = ["real", "integer"]
    n = 1000

    for separator in " \t,":
        extension = {
            " ": "ssv",
            "\t": "tsv",
            ",": "csv",
        }[separator]
        for header in (False, True):
            header_suffix = "header" if header else "noheader"
            filename = f"{name}.{n}.{header_suffix}.{extension}"
            seq = generate_sequence(n)

            with open(filename, "w") as f:
                generate_ctsv(f, separator, header, seq, keys)

    with open(f"{name}.{n}.json", "w") as f:
        json.dump(list(generate_sequence(n)), f, indent=2)

    with open(f"{name}.{n}.jsonl", "w") as f:
        for rec in generate_sequence(n):
            json.dump(rec, f)
            f.write("\n")

    with open(f"{name}.{n}.logfmt", "w") as f:
        seq = generate_sequence(n)
        generate_logfmt(f, seq, keys)
