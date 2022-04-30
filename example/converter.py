import csv
import datetime
import json
import re
import sys
from typing import Dict, Iterator, List, Optional, TextIO

EDT = datetime.timezone(datetime.timedelta(hours=-4), 'EDT')
LONG_BLANKS = re.compile(r'\s+')


def calc_crashed_at(crash_date: str, crash_time: str) -> datetime.datetime:
    # original format is 0:00, force padding zero
    if len(crash_time) == 4:
        crash_time = '0' + crash_time
    return datetime.datetime.strptime(
        crash_date + ' ' + crash_time,
        r'%m/%d/%Y %H:%M'
    ).replace(tzinfo=EDT)


def calc_location(raw_location: str) -> Optional[str]:
    if len(raw_location) == 0:
        return None
    location = raw_location[1:-1]
    # collision on null island!!
    if location == '0.0, 0.0':
        return None
    return location


def sanitize_street_name(raw_street_name: str) -> Optional[str]:
    if len(raw_street_name) == 0:
        return None
    return LONG_BLANKS.sub(' ', raw_street_name).rstrip()


def sanitize_numbers(raw_number: str) -> int:
    if len(raw_number) == 0:
        return 0
    return int(raw_number)


def calc_contributing_factors(row: Dict[str, str]) -> List[str]:
    factors = []
    for num in range(1, 6):
        factor = row[f'CONTRIBUTING FACTOR VEHICLE {num}']
        vehicle = row[f'VEHICLE TYPE CODE {num}']
        if len(factor) and len(vehicle):
            factors.append({
                'vehicle_type': vehicle,
                'factor': factor
            })
    return factors


def convert_row(row: Dict[str, str]) -> Dict[str, str]:
    crashed_at = calc_crashed_at(row['CRASH DATE'], row['CRASH TIME'])
    return {
        'collision_id': row['COLLISION_ID'],
        'crashed_at': crashed_at.isoformat(),
        'borough': row['BOROUGH'] if len(row['BOROUGH']) else None,
        'zipcode': row['ZIP CODE'] if len(row['ZIP CODE']) else None,
        'location': calc_location(row['LOCATION']),
        'on_street_name': sanitize_street_name(row['ON STREET NAME']),
        'cross_street_name': sanitize_street_name(row['CROSS STREET NAME']),
        'off_street_name': sanitize_street_name(row['OFF STREET NAME']),
        'persons_injured': sanitize_numbers(row['NUMBER OF PERSONS INJURED']),
        'persons_killed':  sanitize_numbers(row['NUMBER OF PERSONS KILLED']),
        'pedestrians_injured': sanitize_numbers(row['NUMBER OF PEDESTRIANS INJURED']),
        'pedestrians_killed':  sanitize_numbers(row['NUMBER OF PEDESTRIANS KILLED']),
        'cyclist_injured': sanitize_numbers(row['NUMBER OF CYCLIST INJURED']),
        'cyclist_killed':  sanitize_numbers(row['NUMBER OF CYCLIST KILLED']),
        'motorist_injured': sanitize_numbers(row['NUMBER OF MOTORIST INJURED']),
        'motorist_killed':  sanitize_numbers(row['NUMBER OF MOTORIST KILLED']),
        'contributing_factors': calc_contributing_factors(row)
    }


def convert(source: TextIO) -> Iterator[str]:
    reader = csv.DictReader(source)
    for row in reader:
        doc_row = convert_row(row)
        yield json.dumps(doc_row, ensure_ascii=False)


if __name__ == '__main__':
    source_filepath = sys.argv[1]
    with open(source_filepath) as source:
        for row in convert(source):
            print(row)
