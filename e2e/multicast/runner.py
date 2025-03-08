import random
import os
import time

from types import FunctionType

from mininet.log import setLogLevel, info
from minindn.minindn import Minindn

import multicast_001

def run(scenario: FunctionType, **kwargs) -> None:
    try:
        random.seed(0)

        info(f"===================================================\n")
        start = time.time()
        scenario(ndn, **kwargs)
        info(f'Scenario completed in: {time.time()-start:.2f}s\n')
        info(f"===================================================\n\n")

        # Call all cleanups without stopping the network
        # This ensures we don't recreate the network for each test
        for cleanup in reversed(ndn.cleanups):
            cleanup()
    except Exception as e:
        ndn.stop()
        raise e
    finally:
        # kill everything we started just in case ...
        os.system('pkill -9 ndnd')
        os.system('pkill -9 nfd')

if __name__ == '__main__':
    setLogLevel('info')

    Minindn.cleanUp()
    Minindn.verifyDependencies()
    
    ndn = Minindn(topoFile="../topo.sprint.conf")
    ndn.start()

    run(multicast_001.scenario_old)
    run(multicast_001.scenario_new)

    ndn.stop()
