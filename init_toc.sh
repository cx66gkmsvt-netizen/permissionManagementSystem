#!/bin/bash

# Configuration
API_BASE="http://localhost:8080/api"
USERNAME="admin"
PASSWORD="admin123"
# "To C 事业部" ID is likely 8 based on previous insertion order (7th item start from ID 2)
# Or I could try to fetch it, but let's try 8.
PARENT_ID=8

DEPTS=(
    "教学产品中心"
    "运营中心"
    "用户增长中心"
    "规划中心"
    "服务中心"
    "产研中心"
    "综合品牌中心"
)

# 1. Login
echo "Login as $USERNAME..."
LOGIN_RESP=$(curl -s -X POST "$API_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"userName\": \"$USERNAME\", \"password\": \"$PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESP | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')

if [ -z "$TOKEN" ]; then
    echo "Login failed! Response: $LOGIN_RESP"
    exit 1
fi

echo "Login success."

# 2. Add Departments
SORT=1
for DEPT_NAME in "${DEPTS[@]}"; do
    echo "Adding: $DEPT_NAME to Parent ID $PARENT_ID..."
    
    JSON="{\"parentId\": $PARENT_ID, \"deptName\": \"$DEPT_NAME\", \"sort\": $SORT, \"status\": \"0\"}"
    
    RESP=$(curl -s -X POST "$API_BASE/system/dept" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "$JSON")
        
    echo "Response: $RESP"
    SORT=$((SORT + 1))
done
