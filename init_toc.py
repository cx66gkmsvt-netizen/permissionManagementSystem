import json
import subprocess
import sys

# Configuration
API_BASE = "http://localhost:8080/api"
USERNAME = "admin"
PASSWORD = "admin123"
PARENT_NAME = "To C 事业部"

DEPTS = [
    "教学产品中心",
    "运营中心",
    "用户增长中心",
    "规划中心",
    "服务中心",
    "产研中心",
    "综合品牌中心"
]

def run_curl(method, url, headers=None, data=None):
    cmd = ["curl", "-s", "-X", method, url]
    if headers:
        for k, v in headers.items():
            cmd.extend(["-H", f"{k}: {v}"])
    if data:
        cmd.extend(["-d", json.dumps(data)])
        if "Content-Type" not in (headers or {}):
             cmd.extend(["-H", "Content-Type: application/json"])
    
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception(f"Curl failed: {result.stderr}")
    return result.stdout

def login():
    url = f"{API_BASE}/auth/login"
    data = {"userName": USERNAME, "password": PASSWORD}
    try:
        resp = run_curl("POST", url, headers={"Content-Type": "application/json"}, data=data)
        res_json = json.loads(resp)
        return res_json.get("token")
    except Exception as e:
        print(f"Login failed: {e}")
        sys.exit(1)

def get_dept_id(token):
    url = f"{API_BASE}/system/dept/all" 
    try:
        resp = run_curl("GET", url, headers={"Authorization": f"Bearer {token}"})
        res_json = json.loads(resp)
        depts = res_json.get("data", [])
        if not depts:
             print(f"No data in response: {resp}")
        for dept in depts:
            if dept.get("deptName") == PARENT_NAME:
                return dept.get("deptId")
    except Exception as e:
        print(f"Get departments failed: {e}")
        sys.exit(1)
    return None

def create_dept(token, parent_id, name, sort):
    url = f"{API_BASE}/system/dept"
    data = {
        "parentId": parent_id,
        "deptName": name,
        "sort": sort,
        "status": "0"
    }
    try:
        resp = run_curl("POST", url, headers={"Authorization": f"Bearer {token}", "Content-Type": "application/json"}, data=data)
        print(f"Created {name}: {resp}")
    except Exception as e:
        print(f"Failed to create {name}: {e}")

def main():
    print("Logging in...")
    token = login()
    print("Token acquired.")
    
    print(f"Finding ID for {PARENT_NAME}...")
    parent_id = get_dept_id(token)
    if not parent_id:
        print(f"Department '{PARENT_NAME}' not found!")
        sys.exit(1)
    
    print(f"Found Parent ID: {parent_id}")
    
    sort = 1
    for dept_name in DEPTS:
        print(f"Creating {dept_name}...")
        create_dept(token, parent_id, dept_name, sort)
        sort += 1

if __name__ == "__main__":
    main()
