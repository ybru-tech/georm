# GEOrm

Библиотека была создана для адаптации типов геометрии в [GORM](https://github.com/go-gorm/gorm)

georm изначально создан для работы с PostGIS (PostgreSQL) и может не подойти для работы с другими СУБД.

Основой для географических / геометрических типов является библиотека [go-geom](https://github.com/twpayne/go-geom)

Для обмена геоданными используется *ewkb* формат, серилизация и десерилизация происходит при помощи библиотеки [go-geom/encoding](https://github.com/twpayne/go-geom/tree/master/encoding)

## Features
- Работающая авто-миграция для таблиц с геометрическими типами.
- Возможность создания и получения записей без написания sql, используя только gorm методы.
- Использование бинарного формата в SQL запросах, увеличивает производительность и уменьшает объем трафика
- Метод String, возвращает данные о геометрии в человеко читаемом wkt формате

## Geometry types

- Point
- LineString
- Polygon
- MultiPoint
- MultiLineString
- MultiPolygon
- GeometryCollection

## License

Released under the [MIT Licence](./LICENSE)

Relation: [go-geom](https://github.com/twpayne/go-geom) under the [BSD-2-Clause](https://github.com/twpayne/go-geom/blob/master/LICENSE)
