import re
import socket
from typing import List, Union

from common_ports import ports_and_services

ERROR_INVALID_HOSTNAME = "Error: Invalid hostname"
ERROR_INVALID_IP_ADDRESS = "Error: Invalid IP address"


def get_open_ports(target: str, port_range: List[int], is_verbose: bool = False) -> Union[list, str]:
    if is_ip(target):
        host, address = "", target
        try:
            socket.inet_pton(socket.AF_INET, address)
        except socket.error:
            return ERROR_INVALID_IP_ADDRESS

        try:
            host = socket.gethostbyaddr(address)
            if type(host) is tuple:
                host = host[0]
        except socket.error:
            host = ""
    else:
        host, address = target, ""
        try:
            address = socket.gethostbyname(host)
        except socket.gaierror:
            return ERROR_INVALID_HOSTNAME

    open_ports = []
    start, stop = port_range
    for port in range(start, stop+1):
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(1)
        try:
            s.connect((target, port))
            open_ports.append(port)
        except socket.error:
            pass
        finally:
            s.close()

    if not is_verbose:
        return open_ports

    res = f"Open ports for "
    if host:
        res += f"{host}"
    if address and not host:
        res += f"{address}"
    elif address:
        res += f" ({address})"
    res += "\nPORT     SERVICE\n"
    res += "\n".join(map(lambda p: "{:<9d}{}".format(p, ports_and_services.get(p)), open_ports))
    return res


def is_ip(target: str) -> bool:
    m = re.match(r"^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$", target)
    return bool(m) and all
