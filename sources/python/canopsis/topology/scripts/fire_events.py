from canopsis.topology.manager import TopologyManager
from canopsis.old.rabbitmq import Amqp
from canopsis.engines.core import publish


def fire_events():

    manager = TopologyManager()
    publisher = Amqp()

    graphs = manager.get_graphs()

    for graph in graphs:
        graph.state = 0
        event = graph.get_event(state=0)
        publish(publisher=publisher, event=event)

if __name__ == '__main__':
    fire_events()
