project(
  'herbarium',
  version: '0.3.0',
  meson_version: '>= 1.0.0',
  default_options: [
    'warning_level=2',
  ],
)
i18n = import('i18n')

subdir('po')

app_id = 'ru.ximper.Herbarium'

name_suffix = ''
version_suffix = ''
app_id_suffix = ''

app_id_relevant = app_id + app_id_suffix
app_name = 'Herbarium' + name_suffix

conf = configuration_data()
conf.set_quoted('APP_ID', app_id)
conf.set_quoted('APP_ID_RELEVANT', app_id_relevant)
conf.set_quoted('APP_NAME', app_name)
conf.set_quoted('G_LOG_DOMAIN', app_id_relevant)
conf.set_quoted('VERSION', meson.project_version() + version_suffix)
conf.set_quoted('GETTEXT_PACKAGE', meson.project_name())
conf.set_quoted('DATADIR', get_option('prefix') / get_option('datadir'))
conf.set_quoted('GNOMELOCALEDIR', get_option('prefix') / get_option('localedir'))
conf.set_quoted('HOMEPAGE', 'https://github.com/X1mper/Herbarium')
conf.set_quoted('BUGTRACKER', 'https://github.com/X1mper/Herbarium/issues')
conf.set_quoted('HELP', 'https://t.me/Ximper')

go_root = run_command('printenv', 'GOROOT', check: false)
go_bin = go_root.returncode() != 0 ? find_program('go') : find_program(join_paths(go_root.stdout().strip(), 'bin', 'go'))

# if get_option('shell_completion').enabled()
#   subdir('data/completion')
# endif

custom_target(
  build_by_default: true,
  build_always_stale: true,
  console: true,
  install: true,
  install_dir: get_option('bindir'),
  output: meson.project_name() + '-cli',
  command: [
    'go',
    'build',
    '-trimpath',
    '-ldflags= -s -w',
    '-o', '@OUTPUT@',
    meson.current_source_dir() + '/cli',
  ],
)

custom_target(
  build_by_default: true,
  build_always_stale: true,
  console: true,
  install: true,
  output: meson.project_name() + '-gui',
  install_dir: get_option('bindir'),
  command: [
    'go',
    'build',
    '-trimpath',
    '-ldflags= -s -w',
    '-o', '@OUTPUT@',
    meson.current_source_dir() + '/gui',
  ],
)

subdir('data')
