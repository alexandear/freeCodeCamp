import hashlib


PASSWORDS_FILE = "top-10000-passwords.txt"
SALTS_FILE = "known-salts.txt"
ERROR_NOT_IN_DATABASE = "PASSWORD NOT IN DATABASE"


def crack_sha1_hash(password_hash: str, use_salts: bool = False) -> str:
    if not use_salts:
        if password := crack_password_from_file(PASSWORDS_FILE, password_hash):
            return password
        return ERROR_NOT_IN_DATABASE

    with open(SALTS_FILE, "r") as f:
        while salt := f.readline().rstrip("\n"):
            if password := crack_password_from_file(PASSWORDS_FILE, password_hash, salt):
                return password

    return ERROR_NOT_IN_DATABASE


def crack_password_from_file(file: str, password_hash: str, salt: str = "") -> str:
    def if_hash_match(pwd: str) -> bool:
        return password_hash == hashlib.sha1(pwd.encode()).hexdigest()

    with open(file, "r") as f:
        while password := f.readline().rstrip("\n"):
            if not salt and if_hash_match(password):
                return password

            if if_hash_match(password + salt) or if_hash_match(salt + password):
                return password

    return ""
