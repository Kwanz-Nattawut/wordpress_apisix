# apisix-wp-dashboard

Setup Instructions

Clone the repository

git clone https://github.com/Kwanz-Nattawut/wordpress_apisix.git
cd your-project

Pull docker image and install

docker-compose up -d --build

Step 1: Wordpress

- open your web-browser : localhost:8080 setup the wordpress
- after logined go settings > permalinks > Post name and Save changes ( setting Api wordpress )

Step 2: Apisix-dashboard

- go localhost:9000 - Username : admin , Password : admin
- setting to Upstream and Routes for API endpoint

Step 3: MariaDB

- Create table records and columns ( id , name , value )

Step 4: Go-fiber (back-end)

- main.go for CRUD operations.
