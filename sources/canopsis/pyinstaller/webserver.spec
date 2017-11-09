# -*- mode: python -*-

from PyInstaller.utils.hooks import collect_data_files

block_cipher = None
datas = []
datas.extend(collect_data_files('jsonschema'))

a = Analysis(['../scripts/webserverpy'],
            pathex=['.'],
            hiddenimports=[
                "canopsis",
                "canopsis.engines.dynamic",
                "canopsis.stats.process",
                "canopsis.configuration.driver.file",
                "canopsis.configuration.driver.file.json",
                "canopsis.configuration.driver.file.ini",
                "kombu.transport.pyamqp",
                "canopsis.webcore.services.alerts",
                "canopsis.webcore.services.associativetable",
                "canopsis.webcore.services.auth",
                "canopsis.webcore.services.calendar",
                "canopsis.webcore.services.check",
                "canopsis.webcore.services.context_graph",
                "canopsis.webcore.services.context",
                "canopsis.webcore.services.ctxprop",
                "canopsis.webcore.services.entities",
                "canopsis.webcore.services.event",
                "canopsis.webcore.services.graph",
                "canopsis.webcore.services.gui",
                "canopsis.webcore.services.i18n",
                "canopsis.webcore.services.linklist",
                "canopsis.webcore.services.new_context",
                "canopsis.webcore.services.pbehavior",
                "canopsis.webcore.services.perfdata",
                "canopsis.webcore.services.rest",
                "canopsis.webcore.services.rights",
                "canopsis.webcore.services.session",
                "canopsis.webcore.services.stats",
                "canopsis.webcore.services.storage",
                "canopsis.webcore.services.topology",
                "canopsis.webcore.services.userview",
                "canopsis.webcore.services.vevent",
                "canopsis.webcore.services.watcher",
                "canopsis.webcore.services.weather",
                "canopsis.webcore.services.webcore",
            ],
            hookspath=None,
            runtime_hooks=None,
            cipher=block_cipher,
            datas=datas
)

pyz = PYZ(a.pure, cipher=block_cipher)

exe = EXE(pyz,
          a.scripts,
          a.binaries,
          a.zipfiles,
          name='webserver',
          debug=False,
          strip=False,
          upx=False,
          console=True)

