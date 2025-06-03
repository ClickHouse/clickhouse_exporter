# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).
<!-- insertion marker -->
## [0.1.1](https://github.com/ClickHouse/clickhouse_exporter/releases/tag/0.1.1) - 2025-06-03

<small>[Compare with 0.1.0]()</small>

### Fixed

- style: renamed Exporter to ExporterHolder([f1f7f44](https://ganj-ipe.yaftar.ir/zafir/data_analytics/clickhouse-exporter/commit/f1f7f4425b809c3dc72cc4b8b8161db827ea310c) by smh_tabatabaei).
- style: renamed ParsexxxResponse to ParseResponse([b104e0d](https://ganj-ipe.yaftar.ir/zafir/data_analytics/clickhouse-exporter/commit/b104e0dc9c0ae270c50a4701e4b7fb93a6420445) by smh_tabatabaei)
- fix: bug in table_exporter([d18878e](https://ganj-ipe.yaftar.ir/zafir/data_analytics/clickhouse-exporter/commit/d18878ec1479a8894808262ecf224fb7f1298b84) by smh_tabatabaei)
- feat: added exporter for table metrics([4af960b](https://ganj-ipe.yaftar.ir/zafir/data_analytics/clickhouse-exporter/commit/4af960b28aa50cb61914a2341b751b21a06863fa) by smh_tabatabaei)

<!-- insertion marker -->
## [0.1.0](https://github.com/ClickHouse/clickhouse_exporter/releases/tag/0.1.0) - 2025-06-03

<small>[Compare with first commit](https://github.com/ClickHouse/clickhouse_exporter/compare/3331edb06a30e035617b198db1dffd1145fc56d9...0.1.0)</small>

### Added

- add running check to docker build ([2139fa1](https://github.com/ClickHouse/clickhouse_exporter/commit/2139fa1bba8fa83cd65932a2ca0d5f57812ce2ff) by Slach).
- added disk metrics ([421cf3c](https://github.com/ClickHouse/clickhouse_exporter/commit/421cf3cbc25f28d6ded365a88d8e8700da173acf) by Kelvin).
- Add ClickHouse http client timeout (#27) ([1038e72](https://github.com/ClickHouse/clickhouse_exporter/commit/1038e72f7510b418aa56a12342793edf47b4cb09) by nvartolomei).
- Add badges. ([6b253dc](https://github.com/ClickHouse/clickhouse_exporter/commit/6b253dc4d04b8f5be82b0296a1d1568dc28f6691) by Alexey Palazhchenko).
- Add very basic integration test. ([16b7bac](https://github.com/ClickHouse/clickhouse_exporter/commit/16b7bac5ee7d05eb4f2965167fdca40290c88ad1) by Alexey Palazhchenko).
- add credentials for authorization ([26793e0](https://github.com/ClickHouse/clickhouse_exporter/commit/26793e0cc9a8dc818226603d642cc7d05e1c15d9) by f1yegor).
- Add empty Gopkg.toml to resolve dep complaints ([82ae533](https://github.com/ClickHouse/clickhouse_exporter/commit/82ae533a6969cf6aaa53a3f95a103f5d19bc5858) by Ivan Babrou).
- Add stats for part counts and sizes ([9f859ae](https://github.com/ClickHouse/clickhouse_exporter/commit/9f859aeaf2302bc5ef674797a43acfcb93604074) by Ivan Babrou).
- add asynchronous_metrics table; update prometheus version ([3f7c0f4](https://github.com/ClickHouse/clickhouse_exporter/commit/3f7c0f444d9f29f9e30edf30a2edc97b97aec6d9) by f1yegor).

### Fixed

- fix: now passing all configurations to NewExporter and made its arguments is monadic ([bbb8c38](https://github.com/ClickHouse/clickhouse_exporter/commit/bbb8c3809009d04ab60be6b52a87b56dd2a0c9a1) by smh_tabatabaei).
- fix: moved packages into basic go projects template ([1eac027](https://github.com/ClickHouse/clickhouse_exporter/commit/1eac027e2b602c03a66dfc2ba661d4dbb48b0304) by smh_tabatabaei).
- fix: moved prometheus sample configuration to README.md ([bb268a1](https://github.com/ClickHouse/clickhouse_exporter/commit/bb268a1ff0b68f33092580c88fa7acd11ae518dc) by smh_tabatabaei).
- fix: filtered query result of query_log exporter so it doesn't return data about temporary systematic tables ([8aee53d](https://github.com/ClickHouse/clickhouse_exporter/commit/8aee53d01dac1d8173d008da2698646a88d78abc) by smh_tabatabaei).
- fixe: dockerfile to run the binary go file inside the alpine image ([de84485](https://github.com/ClickHouse/clickhouse_exporter/commit/de84485bddd12ef8e1bcc69761ebb2b9fcbcc4e2) by seyed mohammad hasan tabatabaei).
- fix: go binary wasn't able to run on alpine image now it can ([0fcb2fc](https://github.com/ClickHouse/clickhouse_exporter/commit/0fcb2fc446a47a19a45724ab78748851cf7b4e4f) by smh_tabatabaei).
- fix: removed personal information from sample .env ([f53717f](https://github.com/ClickHouse/clickhouse_exporter/commit/f53717fc20db156af4d6617234852b4579d38cf3) by smh_tabatabaei).
- fix: removed unused files ([c6c3bcd](https://github.com/ClickHouse/clickhouse_exporter/commit/c6c3bcdd5915850d370a4df778dad9c9c393ba9b) by smh_tabatabaei).
- fix: docker-compose now gets variables from .env ([699dfcf](https://github.com/ClickHouse/clickhouse_exporter/commit/699dfcf6fd4af05977b38ab79e7f59e2faec5432) by smh_tabatabaei).
- fix: moved all packages into src ([3922ccc](https://github.com/ClickHouse/clickhouse_exporter/commit/3922ccc8660f65bb60ee2eda2d4c2592c93d5989) by smh_tabatabaei).
- fix: moved exporter queries from exporter.go into specific exporter files ([8d31612](https://github.com/ClickHouse/clickhouse_exporter/commit/8d3161260b1c611d383572dbfd434b5cee0734fd) by smh_tabatabaei).
- fix: returned deleted file .travis.yaml ([5c57e4f](https://github.com/ClickHouse/clickhouse_exporter/commit/5c57e4f4b249503d3743f86b5652b111f7b228e8) by smh_tabatabaei).
- fix: moved disk_metrics to exporters package ([5020fce](https://github.com/ClickHouse/clickhouse_exporter/commit/5020fcee3f8f4908b21d5a7778dfbc1eedaf0f94) by smh_tabatabaei).
- fix: moved parts_metrics to exporters package ([e39a703](https://github.com/ClickHouse/clickhouse_exporter/commit/e39a7033582c2d840e9b58fdc3dba6e20437eac5) by smh_tabatabaei).
- fix: moved event metrics to exporters/event_metrics.go ([cce090d](https://github.com/ClickHouse/clickhouse_exporter/commit/cce090db3f67d79dbf7923fe3407bc5a1d27463a) by smh_tabatabaei).
- fix: query was hardcoded in constructors ([c0ddacb](https://github.com/ClickHouse/clickhouse_exporter/commit/c0ddacbe566732dd102d35fb8907a1fa5e0a9191) by smh_tabatabaei).
- fix: moved async metrics to exporters package ([4089fb8](https://github.com/ClickHouse/clickhouse_exporter/commit/4089fb8e66b63ba76809238992d4a186c9ebce73) by smh_tabatabaei).
- fix: style and polymorphysm ([845a2ef](https://github.com/ClickHouse/clickhouse_exporter/commit/845a2efd63396025724983baf16454519595e0b7) by smh_tabatabaei).
- fix ([e522d22](https://github.com/ClickHouse/clickhouse_exporter/commit/e522d2215c0699f04fd492cce2f27655be2bb7f7) by smh_tabatabaei).
- fix: moved LineResult to util package ([ad30a2e](https://github.com/ClickHouse/clickhouse_exporter/commit/ad30a2e56dbff310e62b55930e9c8a7382cb8552) by smh_tabatabaei).
- fix: moved things a little ([7a285af](https://github.com/ClickHouse/clickhouse_exporter/commit/7a285af49a4d3b6eabfab30effabbdb518214ec1) by smh_tabatabaei).
- fix: moved clickhouse function to util package ([42e5341](https://github.com/ClickHouse/clickhouse_exporter/commit/42e53412ce7ccd82457f9daf7101b4cc61d565be) by smh_tabatabaei).
- fix: moved all packages to ./pkg/ ([6adc3ab](https://github.com/ClickHouse/clickhouse_exporter/commit/6adc3abe956e8bd65e4025c2e697a5c09c4ec6b6) by smh_tabatabaei).
- fix: using regular packge names for import in project ([9a4e743](https://github.com/ClickHouse/clickhouse_exporter/commit/9a4e743e95b786bc4bdf58ae47be66edf0c31898) by smh_tabatabaei).
- fix: changed the dockerfile so it can now be built in the company ([432ef9f](https://github.com/ClickHouse/clickhouse_exporter/commit/432ef9f278ba303fbf4d428536b6d09f61262843) by smh_tabatabaei).
- fix https://github.com/ClickHouse/clickhouse_exporter/issues/84 ([73705de](https://github.com/ClickHouse/clickhouse_exporter/commit/73705de5fe81984ac750145e962273a0e345ad5c) by Slach).
- fix: disk metrics error ([152ced4](https://github.com/ClickHouse/clickhouse_exporter/commit/152ced4323cc6e2cd0bbb9077fa04283a2b7864a) by fuxingZhang).
- fix build.yaml ([8d958b6](https://github.com/ClickHouse/clickhouse_exporter/commit/8d958b6e4dc2a1ff13ce55469894f84f7e11e01b) by Slach).
- fix Dockerfile and Makefile ([403366e](https://github.com/ClickHouse/clickhouse_exporter/commit/403366e05e6c6166ce6cf52b2dab37c8bf095bda) by Slach).
- fix github actions ([7367a98](https://github.com/ClickHouse/clickhouse_exporter/commit/7367a98d54753292bfc4d9eeba7ae896d3b6a10e) by Slach).
- Fixed error: panic: descriptor Desc{fqName: "clickhouse_block_queue_time_dm-2", help: "Number of BlockQueueTime_dm-2 async processed", constLabels: {}, variableLabels: []} is invalid: "clickhouse_block_queue_time_dm-2" is not a valid metric name ([38c5bc6](https://github.com/ClickHouse/clickhouse_exporter/commit/38c5bc6c67be1b68088a88ee6e096e7bc122bb6b) by Mikhail Grigorev).
- fix docker build (#24) ([ba6cbdd](https://github.com/ClickHouse/clickhouse_exporter/commit/ba6cbddc65044f0cbe98daf37bf48471d414b5b0) by Yegor Andreenko).
- fix scrape failures on ch 18.5.x+ (#23) ([ed5fbec](https://github.com/ClickHouse/clickhouse_exporter/commit/ed5fbec15686109cc40fdc6d7fc3d3a98ac49b33) by Roman Tkalenko).
- Fix broken parts scrape. ([d384ca5](https://github.com/ClickHouse/clickhouse_exporter/commit/d384ca546ffaee67e0ebc980bb7d46e32f778ff7) by Alexey Palazhchenko).
- fix Makefile.COMMON: missing curl? ([d6c7673](https://github.com/ClickHouse/clickhouse_exporter/commit/d6c767362a425e1fc773ded05d7a7623e3dd1d75) by f1yegor).
- fix names ([8f99c63](https://github.com/ClickHouse/clickhouse_exporter/commit/8f99c638731702449b3b84d41348fc0514f5646f) by f1yegor).

### Changed

- change url for grafana dashboard ([7b4342d](https://github.com/ClickHouse/clickhouse_exporter/commit/7b4342db0a6c0f2625933ef7db267aa2d65cc149) by Sidorov Pavel).

### Removed

- Remove toInt64 ([a249a90](https://github.com/ClickHouse/clickhouse_exporter/commit/a249a90f89a643bb528c839374df170f911b4e3a) by Mikhail Grigorev).
- remove promu; update docker base ([5781fae](https://github.com/ClickHouse/clickhouse_exporter/commit/5781fae3e7f55b531ae2edcd16611fcff116e832) by yegor).

