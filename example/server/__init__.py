from twirp.asgi import TwirpASGIApp

from ..rpc.haberdasher.service_twirp import HaberdasherServer
from .services import HaberdasherService

service = HaberdasherServer(service=HaberdasherService())
app = TwirpASGIApp()
app.add_service(service)
