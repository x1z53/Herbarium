non_quoted_conf = configuration_data()
non_quoted_conf.set('APP_ID', app_id)
non_quoted_conf.set('APP_ID_RELEVANT', app_id_relevant)
non_quoted_conf.set('HOMEPAGE', 'https://github.com/X1mper/Herbarium')
non_quoted_conf.set('BUGTRACKER', 'https://github.com/X1mper/Herbarium/issues')
non_quoted_conf.set('HELP', 'https://t.me/Ximper')

desktop_file = i18n.merge_file(
  input: configure_file(
    output: app_id_relevant + '.desktop.in',
    input: app_id_relevant + '.desktop.in.in',
    configuration: non_quoted_conf,
  ),
  output: app_id_relevant + '.desktop',
  type: 'desktop',
  po_dir: '../po',
  install: true,
  install_dir: get_option('datadir') / 'applications',
)

desktop_utils = find_program('desktop-file-validate', required: false)
if desktop_utils.found()
  test('Validate desktop file', desktop_utils, args: [desktop_file])
endif

appstream_file = i18n.merge_file(
  input: configure_file(
    output: app_id_relevant + '.metainfo.xml.in',
    input: app_id_relevant + '.metainfo.xml.in.in',
    configuration: non_quoted_conf,
  ),
  output: app_id_relevant + '.metainfo.xml',
  po_dir: '../po',
  install: true,
  install_dir: get_option('datadir') / 'metainfo',
)

appstreamcli = find_program('appstreamcli', required: false)
if appstreamcli.found()
  test(
    'Validate appstream file',
    appstreamcli,
    args: ['validate', appstream_file],
  )
endif

subdir('icons')
