import unittest
import pytest
import _pytest

from src.cosmic import create_app
from src.cosmic.db import delete_db, init_db
import logging
import logging

class TerminalReporter(_pytest.terminal.TerminalReporter):
    def _gettestname(self, rep):
        # actually "rename" original method for clarity
        return super()._getfailureheadline(rep)

    def _getfailureheadline(self, rep):
        # instead of test name
        # (which is originally printed at the top of report)
        # return prefixed name
        return '>>' + self._gettestname(rep)

    def _outrep_summary(self, rep):
        super()._outrep_summary(rep)
        # after printing test report end, print out test name again
        # XXX: here we hard-code color, so red will be used even for passed tests
        # (if passes logging is enabled)
        # You can add some more logic with analyzing report status
        self.write_sep('_', '<<' + self._gettestname(rep), red=True)

@pytest.hookimpl(trylast=True)
def pytest_configure(config):
    # overwrite original TerminalReporter plugin with our subclass
    # we want this hook to be executed after all other implementations
    # to be able to unregister original plugin
    reporter = TerminalReporter(config)
    config.pluginmanager.unregister(name='terminalreporter')
    config.pluginmanager.register(reporter, 'terminalreporter')

class BaseCosmicTest(unittest.TestCase):
    def setUp(self):
        self.app = create_app({
            "TESTING": True,
            "DATABASE": "test_db.sqlite"
        })
        self.client = self.app.test_client()

        with self.app.app_context():
            init_db()
    
    # Remove the standard logger of pytest
    @pytest.fixture(autouse=True)
    def disable_pytest_logging(self, caplog):
        self._caplog = caplog
        caplog.set_level(logging.NOTSET)

    def tearDown(self):
        with self.app.app_context():
            delete_db()


