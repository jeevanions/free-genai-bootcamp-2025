mkdir -p /tmp/nltk_data
chmod -R 777 /tmp/nltk_data
python3 -c "import nltk; nltk.download('stopwords', download_dir='/tmp/nltk_data')"