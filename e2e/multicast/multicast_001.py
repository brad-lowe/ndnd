import random
import os
import time
import sys

from mininet.log import info
from minindn.minindn import Minindn
from minindn.apps.app_manager import AppManager
from minindn.apps.nfd import Nfd

sys.path.append('../')

from fw import NDNd_FW
import dv_util

def scenario_old(ndn: Minindn):
    scenario(ndn, true_multicast=False)

def scenario_new(ndn: Minindn):
    scenario(ndn, true_multicast=True)

def scenario(ndn: Minindn, true_multicast=False, network='/minindn'):
    """
    TBD
    """
    
    fw = NDNd_FW

    info('Starting forwarder on nodes\n')
    AppManager(ndn, ndn.net.hosts, fw)

    dv_util.setup(ndn, network=network)
    dv_util.converge(ndn.net.hosts, network=network)

    # for node in put_nodes:
    #     cmd = f'ndnd put --expose "{network}/{node.name}/test" < {test_file} &'
    #     info(f'{node.name} {cmd}\n')
    #     node.cmd(cmd)

    # info('Waiting for put to complete\n')
    # time.sleep(30)

    # for node in cat_nodes:
    #     put_node = random.choice(put_nodes)
    #     cmd = f'ndnd cat "{network}/{put_node.name}/test" > recv.test.bin 2> cat.log'
    #     info(f'{node.name} {cmd}\n')
    #     node.cmd(cmd)
    #     if node.cmd(f'diff {test_file} recv.test.bin').strip():
    #         info(node.cmd(f'cat cat.log'))
    #         raise Exception(f'Test file contents do not match on {node.name}')
