DIR="$PWD/db/data"

#required subdirectories to fetch data only
PATHS=("talents" "talentmaterialtypes" "weapons" "weaponmaterialtypes" "artifacts" "domains")

if [ -d "$DIR" ]; then
    if [ $PWD != $DIR ]; then
        cd db/data
    fi
    echo "Fetching data files to ${DIR}..."
else
    echo "Initializing data folder to ${DIR}..."
    mkdir db/data && cd db/data
    git init
    # use sparse checkout to fetch only required subdirectories
    git config core.sparseCheckout true
    git remote add db "https://github.com/theBowja/genshin-db.git"
    for path in "${PATHS[@]}"; do
        echo "Fetching ${path}..."
        echo "src/data/English/${path}" >>.git/info/sparse-checkout
    done
fi

git fetch --depth=1 db main
git checkout main
