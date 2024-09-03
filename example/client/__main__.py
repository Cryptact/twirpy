from twirp.context import Context
from twirp.exceptions import TwirpServerException

from ..rpc.haberdasher import service_pb2 as pb
from ..rpc.haberdasher.service_twirp import HaberdasherClient


def main():
    client = HaberdasherClient("http://localhost:8080")
    try:
        response = client.MakeHat(
            ctx=Context(),
            request=pb.Size(inches=12),
        )
        print(f"I have a nice new hat:\n{response}")
    except TwirpServerException as e:
        print(e.code, e.message, e.meta, e.to_dict())


if __name__ == "__main__":
    main()
