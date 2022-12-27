from pytest import fixture

def pytest_addoption(parser):
    parser.addoption("--url", action="store")

@fixture()
def name(request):
    return request.config.getoption("--url")
