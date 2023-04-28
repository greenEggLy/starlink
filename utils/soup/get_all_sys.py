import requests

# from bs4 import BeautifulSoup, NavigableString, Tag
import bs4
import re
import mysql.connector


def write_to_database(tables):
    cnx = connection.MySQLConnection(
        user="root", password="", host="127.0.0.1", database="starlink"
    )
    add_query = (
        "INSERT INTO systems(name)" " VALUES(%s)" " ON DUPLICATE KEY UPDATE name = %s"
    )
    cursor = cnx.cursor()
    cursor.executemany(add_query, tables)
    cnx.commit()
    cursor.close()
    cnx.close()


def fetch_file():
    html = requests.get("https://celestrak.org/NORAD/elements/").text
    soup = bs4.BeautifulSoup(html, "lxml")

    hFile = open("utils/soup/hfile.html", "w")
    hFile.write(str(soup))
    hFile.close()


def read_file():
    hFile = open("utils/soup/hfile.html", "r")
    html_text = hFile.read()
    tables = []

    tag_reg = re.compile(r'FORMAT=tle" title="TLE Data">(([\s\S])*?)</a>')
    tags = tag_reg.findall(html_text)
    for tag in tags:
        tables.append(tag[0])
    return tables


def main():
    fetch_file()
    tables = read_file()
    # print(len(tables))
    # write_to_database(tables)
    print(tables)


if __name__ == "__main__":
    main()

# innertable = list(soup.table.children)
# table_body = innertable[0]

# tables = table_body.find_all("tbody")
# for table in tables:
#     rows = table.find_all("tr")
#     for row in rows:
#         if isinstance(row, NavigableString):
#             continue
#         cols = row.find_all("td")
#         for col in cols:
#             if "A" <= row.text[0] <= "Z" or "a" <= row.text[0] <= "z":
#                 print("table name:" + row.text)
