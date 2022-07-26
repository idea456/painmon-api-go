DIR="$PWD/db/data"

if [ -d "$DIR" ]; then
    echo "Fetching data files to ${DIR}..."
else
    echo "Initializing data folder to ${DIR}..."
    mkdir db/data && cd db/data
    git init
    git remote add -f origin "https://github.com/theBowja/genshin-db.git"
    git config core.sparseCheckout true
    echo "src/data/English" >>.git/info/sparse-checkout
fi

git pull origin main
