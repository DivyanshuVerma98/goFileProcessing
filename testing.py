import csv
import sys
import copy
import pandas as pd
import random

try:
    number_of_rows = int(sys.argv[1])
except (IndexError, ValueError):
    number_of_rows = 100000

appro_user = [
    [
        "TPOSP00615",
        "BASKARAN M",
        "SANJITH100",
        "07/08/2023",
        "HDFC ERGO",
        "SUGANYA M",
        "Motor",
        "2W",
        "SCOOTER",
        1,
        "Comp",
        "550",
        "465.62",
        "465.62",
        "",
        "",
        "TN65AJ3607",
        "TN65",
        "TAMIL NADU",
        "rto cluster",
        "RAMANATHAPURAM",
        "",
        "PETROL",
        "cpa",
        "125",
        "NA",
        "NO",
        "2",
        "2019",
        "",
        "Suzuki",
        "Access",
        "CITIN23377280784",
        "07/08/2023",
        "1007896",
        "T+7_Slot1",
        "Net",
        "28%",
        "128.046",
        "Net",
        "28%",
        "128.0455",
        "ZIBPL",
        "",
        "approved",
        "",
        ""
    ]

]

enricher_user = [
    [
        "TPOSP00615",
        "BASKARAN M",
        "SANJITH100",
        "07/08/2023",
        "HDFC ERGO",
        "SUGANYA M",
        "Motor",
        "2W",
        "SCOOTER",
        1,
        "Comp",
        "550",
        "465.62",
        "465.62",
        "",
        "",
        "TN65AJ3607",
        "TN65",
        "TAMIL NADU",
        "rto cluster",
        "RAMANATHAPURAM",
        "",
        "PETROL",
        "cpa",
        "125",
        "NA",
        "NO",
        "2",
        "2019",
        "",
        "Suzuki",
        "Access",
        "CITIN23377280784",
        "07/08/2023",
        "1007896",
        "T+7_Slot1",
        "Net",
        "28%",
        "128.046",
        "Net",
        "28%",
        "128.0455",
        "ZIBPL",
        "",
        "SENT TO NEXT STAGE",
        "",
        ""
    ]

]

maker_user = [
    [
        "Primary",
        "TPOSP00615",
        "BASKARAN M",
        "SANJITH100",
        "07/08/2023",
        "HDFC ERGO",
        "SUGANYA M",
        "Motor",
        "2W",
        "SCOOTER",
        1,
        "Comp",
        "550",
        "465.62",
        "465.62",
        "",
        "",
        "TN65AJ3607",
        "TN65",
        "TAMIL NADU",
        "rto cluster",
        "RAMANATHAPURAM",
        "",
        "PETROL",
        "cpa",
        "125", # cc
        "NA",
        "NO",
        "2",
        "2019",
        "",
        "Suzuki",
        "Access",
        "CITIN23377280784",
        "07/08/2023",
        "1007896",
        "T+7_Slot1",
        "Net",
        "28%",
        "128.046",
        "Net",
        "28%",
        "128.0455",
        "ZIBPL",
        ""
    ],
    [
        "Primary",
        "TPOSP00615",
        "BASKARAN M",
        "SANJITH100",
        "07-08-2023",
        "HDFC ERGO",
        "SUGANYA M",
        "Motor",
        "2W",
        "SCOOTER",
        1,
        "Comp",
        "550",
        "465.62",
        "465.62",
        "",
        "",
        "TN65AJ3607",
        "TN65",
        "TAMIL NADU",
        "rto cluster",
        "RAMANATHAPURAM",
        "",
        "PETROL",
        "cpa",
        "125", # cc
        "NA",
        "NO",
        "2",
        "2019",
        "",
        "Suzuki",
        "Access",
        "CITIN23377280784",
        "07/08/2023",
        "1007896",
        "T+7_Slot1",
        "Net",
        "28%",
        "128.046",
        "Net",
        "28%",
        "128.0455",
        "ZIBPL",
        ""
    ]

]

data = []
name_format = "NAME{i}"
code_format = "PCODE{i:020d}"
number_of_rows = 1000000
count = 1
for i in range(number_of_rows):
    # ll = random.choice(appro_user)
    # ll = appro_user[0]
    # ll = enricher_user[0]
    ll = random.choice(maker_user)
    aa = copy.deepcopy(ll)
    aa[10]=count
    count+=1
    data.append(aa)

OTHER_HEADERS = ["RM code", "RM Name1", "Child ID", "Booking Date(Click to select Date)", 
           "Insurer Name", "Insured Name", "Major Categorisation( Motor/Life/ Health)", "Product", 
           "Product type", "Policy number", "Plan type", "Premium", "Net premium", "OD", "TP", 
           "Commissionable premium", "Registrationno", "RTO Code", "State", "RTO Cluster", "City",
            "Insurer_biff", "Fuel type", "CPA", "CC", "GVW", "NCB Type", "Seating Capacity", 
            "Vehicle registration year", "Discount %", "Make", "Model", "UTR", "UTR Date", "UTR Amount", 
            "Slot (Payment batch)", "Paid On IN", "IN%", "IN Amount", "Paid On Out", "Out%", "Out Amount", 
            "Co type ( ZIBPL/I CARE/Individual)", "Remarks", "Status", "User Remark", "Adjustment Amount"]
MAKER_HEADERS = ["Transaction Type", "RM code", "RM Name1", "Child ID", "Booking Date(Click to select Date)", 
           "Insurer Name", "Insured Name", "Major Categorisation( Motor/Life/ Health)", "Product", 
           "Product type", "Policy number", "Plan type", "Premium", "Net premium", "OD", "TP", 
           "Commissionable premium", "Registrationno", "RTO Code", "State", "RTO Cluster", "City",
            "Insurer_biff", "Fuel type", "CPA", "CC", "GVW", "NCB Type", "Seating Capacity", 
            "Vehicle registration year", "Discount %", "Make", "Model", "UTR", "UTR Date", "UTR Amount", 
            "Slot (Payment batch)", "Paid On IN", "IN%", "IN Amount", "Paid On Out", "Out%", "Out Amount", 
            "Co type ( ZIBPL/I CARE/Individual)", "Remarks"]
with open(f"test-data-set-{number_of_rows}.csv", "w", newline="") as csvfile:
    filewriter = csv.writer(csvfile, delimiter=",", quotechar='"')
    filewriter.writerows([MAKER_HEADERS])
    filewriter.writerows(data)

# def append_to_csv(file_path, data):
#     # Create a DataFrame with the new data
#     df = pd.DataFrame(data)

#     # Append the DataFrame to the CSV file
#     with open(file_path, "a", newline="") as csvfile:
#         df.to_csv(csvfile, header=False, index=False)

# # Example usage:
# file_path = f"test-data-set-{number_of_rows}.csv"

# # Append data to the CSV file
# append_to_csv(file_path, data)