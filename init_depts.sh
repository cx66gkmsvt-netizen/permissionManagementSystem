#!/bin/bash

# Configuration
API_BASE="http://localhost:8080/api"
USERNAME="admin"
PASSWORD="admin123"

# Departments to add (Parent ID = 1, based on 'Main Company')
# Assuming parentId 1 exists (created by initData as '总公司')
PARENT_ID=1

DEPTS=(
    "总经办"
    "行政部"
    "人力资源部"
    "财务部"
    "技术部"
    "To B 事业部"
    "To C 事业部"
    "深圳校区"
    "香港校区"
    "兼职"
    "其他"
    "客户对接专用"
    "其他（待设置部门）"
)

# 1. Login to get token
echo "Login as $USERNAME..."
LOGIN_RESP=$(curl -s -X POST "$API_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"userName\": \"$USERNAME\", \"password\": \"$PASSWORD\"}")

# Extract token using grep/sed (simple parser)
TOKEN=$(echo $LOGIN_RESP | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')

if [ -z "$TOKEN" ]; then
    echo "Login failed! Response: $LOGIN_RESP"
    exit 1
fi

echo "Login success. Token acquired."

# 2. Add Departments
SORT=1
for DEPT_NAME in "${DEPTS[@]}"; do
    echo "Adding department: $DEPT_NAME (Sort: $SORT)..."
    
    # Construct JSON using manual string building to avoid jq dependency
    JSON="{\"parentId\": $PARENT_ID, \"deptName\": \"$DEPT_NAME\", \"sort\": $SORT, \"status\": \"0\"}"
    
    RESP=$(curl -s -X POST "$API_BASE/system/dept" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "$JSON")
        
    echo "Response: $RESP"
    
    SORT=$((SORT + 1))
done

echo "All departments processed."
