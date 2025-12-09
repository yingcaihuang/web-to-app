package com.webtoapp.core.apkbuilder

import android.content.Context
import java.io.File
import java.text.SimpleDateFormat
import java.util.*

/**
 * 构建日志记录器 - 将日志写入文件方便调试
 */
object BuildLogger {
    
    private var logFile: File? = null
    private val dateFormat = SimpleDateFormat("HH:mm:ss.SSS", Locale.getDefault())
    
    fun init(context: Context, appName: String) {
        val logDir = File(context.getExternalFilesDir(null), "build_logs").apply { mkdirs() }
        val timestamp = SimpleDateFormat("yyyyMMdd_HHmmss", Locale.getDefault()).format(Date())
        logFile = File(logDir, "build_${timestamp}_${appName}.log")
        logFile?.writeText("===== 构建日志开始: $appName =====\n")
        log("日志文件: ${logFile?.absolutePath}")
    }
    
    fun log(message: String) {
        val time = dateFormat.format(Date())
        val line = "[$time] $message\n"
        logFile?.appendText(line)
        android.util.Log.d("ApkBuilder", message)
    }
    
    fun getLogFilePath(): String? = logFile?.absolutePath
    
    fun getLogContent(): String = logFile?.readText() ?: "无日志"
}
