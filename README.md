# Email Spam Classifier in Go

A simple email spam classifier written in Go using a **Bag of Words (BoW)** approach. It analyzes word frequencies in emails and uses logarithmic probability scores to classify emails as spam or ham (legitimate email).

## How It Works

- **Training:** Builds separate word-frequency maps for spam and ham emails using the Enron dataset (`enron1`–`enron5`).
- **Tokenization:** Splits email content into words and converts them to uppercase.
- **Word Filtering:** Ignores words whose combined frequency in the training data is below a threshold of 100.
- **Classification:** Calculates logarithmic scores using word frequencies and class probabilities, then predicts whether each email is spam or ham.
- **Evaluation:** Tests the classifier on `enron6`, keeping the evaluation data separate from the training data.

## Setup

Extract the dataset:

```bash
tar -xf data/data.tar
```

## Run

```bash
go run main.go
```

## Sample Output

```text
Training...
Clasifying ham...
Spamcount = 331
 HamCount = 1169
Clasifying spam...
Spamcount = 4206
 HamCount = 294
```

The first set of counts shows predictions for emails in the ham directory, while the second shows predictions for emails in the spam directory. Comparing these counts helps evaluate how well the classifier distinguishes spam from legitimate emails.

## Reference

[YouTube Tutorial](https://www.youtube.com/watch?v=JsfOXk7qmSM&t=5s)

## Concepts Practiced

- Go file handling and directory traversal
- Maps, slices, and string processing
- Word-frequency analysis
- Logarithmic probability calculations
- Basic text classification and evaluation

