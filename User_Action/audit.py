with open('audit.txt', 'r') as audit_file:
    paths = audit_file.readlines()

with open('pass.txt', 'a') as pass_file, open('audit.txt', 'w') as audit_file:
    for path in paths:
        print(f"Do you want to add {path.strip()} to pass.txt? (y/n)")
        answer = input()
        if answer == 'y':
            pass_file.write(path)
        else:
            audit_file.write(path)

print("Done!")
