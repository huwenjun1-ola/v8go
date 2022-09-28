#ifdef COMPILE_EXPORT
#ifndef INSPECTOR_V8INSPECTORCLIENTIMPL_H
#define INSPECTOR_V8INSPECTORCLIENTIMPL_H

#include <locale>
#include <codecvt>
#include "v8-inspector.h"
#include "v8.h"

class V8InspectorChannelImpl : public v8_inspector::V8Inspector::Channel {
public:
    V8InspectorChannelImpl(const std::unique_ptr<v8_inspector::V8Inspector> &InV8Inspector, const int32_t InCxtGroupID,
                           int32_t InInspectorId) {
        InspectorId = InInspectorId;
        v8_inspector::StringView DummyState;
        V8InspectorSession = InV8Inspector->connect(InCxtGroupID, this, DummyState);
    }


    virtual ~V8InspectorChannelImpl() override {
        OnSendMessage = nullptr;
        V8InspectorSession.reset();
    }

    void DispatchProtocolMessage(const std::string &Message) {
        const auto MessagePtr = reinterpret_cast<const uint8_t *>(Message.c_str());
        const auto MessageLen = (size_t) Message.length();

        v8_inspector::StringView StringView(MessagePtr, MessageLen);
        V8InspectorSession->dispatchProtocolMessage(StringView);
    }

private:
    void SendMessage(v8_inspector::StringBuffer &MessageBuffer) {
        v8_inspector::StringView MessageView = MessageBuffer.string();

        std::string Message;
        if (MessageView.is8Bit()) {
            Message = reinterpret_cast<const char *>(MessageView.characters8());
        } else {
#if PLATFORM_WINDOWS
            std::wstring_convert<std::codecvt_utf8_utf16<uint16_t>, uint16_t> Conv;
            const uint16_t* Start = MessageView.characters16();
#else
            std::wstring_convert<std::codecvt_utf8_utf16<char16_t>, char16_t> Conv;
            const char16_t *Start = reinterpret_cast<const char16_t *>(MessageView.characters16());
#endif
            Message = Conv.to_bytes(Start, Start + MessageView.length());
        }

        if (OnSendMessage)
            OnSendMessage(InspectorId, Message.c_str());
    }

    void sendResponse(int CallID, std::unique_ptr<v8_inspector::StringBuffer> Message) override {
        SendMessage(*Message);
    }

    void sendNotification(std::unique_ptr<v8_inspector::StringBuffer> Message) override {
        SendMessage(*Message);
    }

    void flushProtocolNotifications() override {

    }

    std::unique_ptr<v8_inspector::V8InspectorSession> V8InspectorSession;

public:
    std::function<void(int32_t, const char *)> OnSendMessage;

    int32_t InspectorId;
};

class V8InspectorClientImpl : public v8_inspector::V8InspectorClient {
public:


    V8InspectorClientImpl(v8::Local<v8::Context> InContext, int32_t inspectorId) {
        Isolate = InContext->GetIsolate();
        Context.Reset(Isolate, InContext);
        IsAlive = false;
        CtxGroupID = 1;
        const uint8_t CtxNameConst[] = "V8InspectorContext";
        v8_inspector::StringView CtxName(CtxNameConst, sizeof(CtxNameConst) - 1);
        V8Inspector = v8_inspector::V8Inspector::create(Isolate, this);
        V8Inspector->contextCreated(v8_inspector::V8ContextInfo(InContext, CtxGroupID, CtxName));
        V8InspectorChannel = std::make_unique<V8InspectorChannelImpl>(V8Inspector, CtxGroupID, inspectorId);
        InspectorId = inspectorId;
        IsAlive = true;
        IsPaused = false;
    }

    virtual ~V8InspectorClientImpl() {
        Close();
    }


    void runMessageLoopOnPause(int ContextGroupId) override {
        if (IsPaused) {
            return;
        }
        IsPaused = true;
        while (IsPaused) {
            if (TickFunc != nullptr) {
                TickFunc(InspectorId);
            }
        }
    }

    void runIfWaitingForDebugger(int contextGroupId) override {
    }

    void quitMessageLoopOnPause() override {
        IsPaused = false;
    }


    v8::Isolate *Isolate;

    v8::Persistent<v8::Context> Context;

    std::unique_ptr<v8_inspector::V8Inspector> V8Inspector;

    int32_t CtxGroupID;


    bool IsAlive;

    bool IsPaused;

    void Close() {
        if (IsAlive) {
            V8InspectorChannel.reset();
            v8::Isolate::Scope IsolateScope(Isolate);
            v8::HandleScope HandleScope(Isolate);
            V8Inspector->contextDestroyed(Context.Get(Isolate));
            IsAlive = false;
            IsPaused = false;
        }
    }

    void OnReceiveMessage(char *Message) {
        V8InspectorChannel->DispatchProtocolMessage(Message);
    }


    std::unique_ptr<V8InspectorChannelImpl> V8InspectorChannel;

    void (*TickFunc)(int32_t);

    int32_t InspectorId;
};

bool bindMessageSendFuncToClient(void *clientPtr, void (*in)(int32_t, const char *)) {
    V8InspectorClientImpl *client = (V8InspectorClientImpl *) (clientPtr);
    client->V8InspectorChannel->OnSendMessage = in;
    return false;
}

bool bindTickFuncToClient(void *clientPtr, void (*in)(int32_t)) {
    V8InspectorClientImpl *client = (V8InspectorClientImpl *) (clientPtr);
    client->TickFunc = in;
    return false;
}


void onReceiveMessage(void *clientPtr, char *Message) {
    V8InspectorClientImpl *client = (V8InspectorClientImpl *) (clientPtr);
    v8::Isolate *iso = client->Isolate;
    v8::Locker locker(iso);
    v8::Isolate::Scope isolate_scope(iso);
    v8::HandleScope handle_scope(iso);
    std::string data(Message);
    client->V8InspectorChannel->DispatchProtocolMessage(data);
}

#endif
#endif //INSPECTOR_V8INSPECTORCLIENTIMPL_H
