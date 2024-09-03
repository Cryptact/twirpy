import random

from twirp.context import Context
from twirp.exceptions import InvalidArgument

from ..rpc.haberdasher import service_pb2 as pb


class HaberdasherService(object):
    def MakeHat(self, context: Context, size: pb.Size) -> pb.Hat:
        if size.inches <= 0:
            raise InvalidArgument(argument="inches", error="I can't make a hat that small!")
        return pb.Hat(
            size=size.inches,
            color=random.choice(["white", "black", "brown", "red", "blue"]),
            name=random.choice(["bowler", "baseball cap", "top hat", "derby"])
        )
